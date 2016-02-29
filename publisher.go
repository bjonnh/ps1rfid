package main

import (
	"fmt"
	"os"
	zmq "github.com/pebbe/zmq4"
)

type Publisher interface {
	SendMessage(string, string) bool
}

type ZMQPublisher struct {
	publisher *zmq.Socket
}

func NewZMQPublisher() *ZMQPublisher {
	publisher := ZMQPublisher{}
	//Configure ZMQ publisher
	_publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		fmt.Println("Error creating the ZMQ publisher")
		fmt.Println(err)
		os.Exit(1)
	}
	publisher.publisher = _publisher
	publisher.publisher.Bind("tcp://*:5556")
	return &publisher
}

func (publisher *ZMQPublisher) SendMessage(id string, msg string) bool {
	publisher.publisher.SendMessage(id, msg)
	return true
}


type FakePublisher struct {
	publisher bool
}

func NewFakePublisher() *FakePublisher {
	publisher := FakePublisher{}
	fmt.Println("FAKEPUBLISHER: Initializing")
	return &publisher
}

func (publisher *FakePublisher) SendMessage(id string, msg string) bool {
	fmt.Printf("FAKEPUBLISHER: I received id:%s msg:%s\n",id, msg)
	return true
}
