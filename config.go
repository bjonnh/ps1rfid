package main

var cacheDB CacheDB
var board Board
var source Source
var publisher Publisher
var auth Auth

func config() {

	var db_filename string = "rfid-tags.db"
	var bucket_name string = "RFIDBucket"
	var ps1auth_url string = "https://members.pumpingstationone.org/rfid/check/FrontDoor/"
	var serial_port string = "/dev/ttyUSB0"
	var serial_port_speed int = 9600	
	
	// Open or create the cache db
	// And get an object that allows to play in the bucket
	cacheDB = NewBoltDb(db_filename, bucket_name)
	board = NewBeagleBone("beaglebone", "splate", "P9_11")
	source = NewSerialSource(serial_port, serial_port_speed)
	publisher = NewZMQPublisher()
	auth = NewPS1Auth(ps1auth_url)
	
	// This is a test behavior with all functions being in debug mode
	
	//board = NewFakeBoard() 
	//source = NewFakeSource()  // The fake source sends the code YOU'REAFAKE once, then it exits at the second call
	//publisher = NewFakePublisher() // Writes to STDOUT
	//auth = NewFakeAuth(auth_response{code:0,msg:"RFID Accepted"}) // Sends the auth_response you want
	//auth := NewFakeAuth(auth_response{code:1,msg:"RFID Failed"}) 
	//auth := NewFakeAuth(auth_response{code:2,msg:"RFID Not found"})
	//auth := NewFakeAuth(auth_response{code:3,msg:"Auth system error"})
}
