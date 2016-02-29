package main

import (
	"fmt"
	"net/http"
)

var db_filename string = "rfid-tags.db"
var bucket_name string = "RFIDBucket"
var ps1auth_url = "https://members.pumpingstationone.org/rfid/check/FrontDoor/"

func main() {
	var code string
	// Now open the cache db to check if it's already here
	// And get an object that allows to play in the bucket
	cacheDB := NewCacheDb(db_filename, bucket_name)

	// This is for a normal behavior
	
	board := NewBeagleBone("beaglebone", "splate", "P9_11")
	source := NewSerialSource("/dev/ttyUSB0", 9600)
	publisher := NewZMQPublisher()
	auth := NewPS1Auth(ps1auth_url)

	// This is a test behavior with all functions being in debug mode

	//board := NewFakeBoard() 
	//source := NewFakeSource()  // The fake source sends the code YOU'REAFAKE once, then it exits at the second call
	//publisher := NewFakePublisher() // Writes to STDOUT
	//auth := NewFakeAuth(auth_response{code:0,msg:"RFID Accepted"}) // Sends the auth_response you want
	//auth := NewFakeAuth(auth_response{code:1,msg:"RFID Failed"}) 
	//auth := NewFakeAuth(auth_response{code:2,msg:"RFID Not found"})
	//auth := NewFakeAuth(auth_response{code:3,msg:"Auth system error"})

	// Tell the board system to use this publisher
	board.setPublisher(publisher)
	// Configure the local API (/ for the last code, /open and /close)
	
	go http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Code: "))
		w.Write([]byte(code))
	})
	// the anonymous function here allows us to call openDoor with splate remaining in scope
	go http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
		board.openDoor()
	})
	
	go http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
		board.closeDoor()
	})

	go http.ListenAndServe(":8080", nil)

	var resp auth_response
	for {
		code = source.getFullCode()				
		// Before checking the site for the code, let's check our cache
		if cacheDB.checkCacheDBForTag(code) == false {
			resp = auth.request(code)
			fmt.Printf("I received response: %s\n",resp)
			if resp.code == 0 {
				// We got 200 back, so we're good to add this
				// tag to the cache
				cacheDB.addTagToCacheDB(code)
				
				fmt.Println("Success!")
				publisher.SendMessage("door.rfid.accept", resp.msg)
				board.openDoor()
			} else if resp.code == 1 {
				fmt.Println("Membership status: Expired")
				publisher.SendMessage("door.rfid.deny", resp.msg)
			} else if resp.code == 2 {
				fmt.Println("Code not found")
				publisher.SendMessage("door.rfid.deny", resp.msg)
			} else if resp.code == 3 {
				fmt.Println("Auth server error")
				publisher.SendMessage("door.rfid.error", resp.msg)
			} else {
				fmt.Println("Unknown Auth error")
				publisher.SendMessage("door.rfid.snafu", resp.msg)
			}
  		} else {
			// If we're here, we found the tag in the cache, so
			// let's just go and open the door for 'em
			fmt.Println("Success!")
			publisher.SendMessage("door.rfid.accept", "RFID Accepted")
			board.openDoor()
		}	
	}
	cacheDB.Close()
}
