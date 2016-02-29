package main

import (
	"fmt"
	"os"
	"io"
	"time"
	"github.com/tarm/goserial"
)

type Source interface {
	getFullCode() string
}

type SerialSource struct {
	source io.ReadWriteCloser
}

func NewSerialSource(dev string, speed int) *SerialSource {
	source := SerialSource{}
	c := &serial.Config{Name: dev, Baud: speed}
	u, err := serial.OpenPort(c)
	if err != nil {
		fmt.Printf("Impossible to open port %v (speed=%v)\n", dev, speed)
		fmt.Println(err)
		os.Exit(1)
	}
	source.source = u
	return &source
}

func (p *SerialSource) getFullCode() string {
	buf := make([]byte, 16)
	source := *p
	n, err := io.ReadFull(source.source, buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// We need to strip the stop and start bytes from the tag
	// so we only assign a certain range of the slice
	return string(buf[1 : n-3])
}



type FakeSource struct {
	count bool
}

func NewFakeSource() *FakeSource {
	source := FakeSource{}
	fmt.Println("FAKESOURCE: Initializing")
	source.count = false
	return &source
}

func (source *FakeSource) getFullCode() string {
	if source.count {
		fmt.Println("FAKESOURCE: I'm done with you, but lets wait some seconds")

		time.Sleep(time.Second * 10)
		os.Exit(1)
	} else {
		fmt.Println("FAKESOURCE: Sending the code")
		source.count = true
	}
	return "YOU'REAFAKE"
}
