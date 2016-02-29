// Defined auth responses

// 0: Authorized
// 1: Denied
// 2: Not found
// 3: Auth system error

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

type Auth interface {
	request(string) (bool,error)
}

type auth_response struct {
	code int
	msg string
}

type PS1Auth struct {
	url string
}

func NewPS1Auth(url string) *PS1Auth {
	auth := PS1Auth{url: url}
	return &auth
}

func (auth *PS1Auth) request(code string) auth_response {
	var request bytes.Buffer
	request.WriteString(auth.url)
	request.WriteString(code)
	resp, err := http.Get(request.String())
	if err != nil {
		fmt.Println("Whoops!")
		return auth_response{code:3, msg: fmt.Sprintf("Auth system error %s",err)}
		os.Exit(1)
	}
	var response auth_response
	if (resp.StatusCode == 200) {
		response = auth_response{code:0,msg:"RFID Accepted"}
	} else if resp.StatusCode == 403 {
		response = auth_response{code:1,msg:"RFID Denied"}
	} else {
		response = auth_response{code:2,msg:"RFID not found"}
	}
	return response
}


type FakeAuth struct {
	response auth_response
}

func NewFakeAuth(response auth_response) *FakeAuth {
	auth := FakeAuth{response: response}
	fmt.Printf("FAKEAUTH: Initializing with response: %s\n", response)
	return &auth
}

func (auth *FakeAuth) request(code string) auth_response {
	fmt.Println("FAKEAUTH: I received code %s",code)
	return auth.response
}
