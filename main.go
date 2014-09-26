package main

import (
	"bytes"
	"net/http"
	"os"
)

import "fmt"

func main() {
	var code string
	for {
		fmt.Print("Enter code: ")
		fmt.Scanf("%s", &code)
		var request bytes.Buffer
		request.WriteString("https://members.pumpingstationone.org/rfid/check/FrontDoor/")
		request.WriteString(code)
		resp, err := http.Get(request.String())
		if err != nil {
			fmt.Printf("Whoops!")
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			fmt.Println("Success!")
		} else if resp.StatusCode == 403 {
			fmt.Println("Membership status: Expired")
		} else {
			fmt.Println("Code not found")
		}
	}

}
