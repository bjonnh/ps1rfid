package main

import (
	"fmt"
	"net/http"
)

func main() {
	var code string
	config()
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
