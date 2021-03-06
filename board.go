package main

import (
	"fmt"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot"
	"os"
	"time"
)

type Board interface {
	openDoor() bool
	closeDoor() bool
	setPublisher(Publisher) bool
}

type BeagleBoneBlack struct {
	publisher Publisher
	pin gpio.DirectPinDriver
	board *beaglebone.BeagleboneAdaptor
	splate *gpio.DirectPinDriver
}

func NewBeagleBone(name string, pinname string, pin string) *BeagleBoneBlack {
	board := BeagleBoneBlack{}
	// gobot as clearly an issue, if you don't run it on the beaglebone,
	// it just crashes stupidly. This is an error that should be catched
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("There is no beaglebone connected, I'll just exit silently and you will try to fix that ok?")
			os.Exit(42)
		}
	}()
	board.board = beaglebone.NewBeagleboneAdaptor("beaglebone")
	//NewDirectPinDriver returns a pointer - this wasn't immediately obvious to me
	board.splate = gpio.NewDirectPinDriver(board.board, "splate", "P9_11")
	return &board
}

func (board *BeagleBoneBlack) openDoor() bool {
	//board.pin.DigitalWrite(1)
	board.publisher.SendMessage("door.state.unlock", "Door Unlocked")
	gobot.After(5*time.Second, func() {
		board.closeDoor()
	})
	return true
}

func (board *BeagleBoneBlack) closeDoor() bool {
	//board.pin.DigitalWrite(0)
	board.publisher.SendMessage("door.state.lock", "Door Locked")
	return true
}

func (board *BeagleBoneBlack) setPublisher(publisher Publisher) bool {
	board.publisher = publisher
	return true
}

type FakeBoard struct {
	publisher Publisher
}

func NewFakeBoard() *FakeBoard {
	var board FakeBoard
	fmt.Println("FAKEBOARD: Initializing")
	board = FakeBoard{}
	return &board
}

func (board *FakeBoard) openDoor() bool {
	fmt.Println("FAKEBOARD: Open door")
	board.publisher.SendMessage("door.state.unlock", "Door Unlocked")
	gobot.After(5*time.Second, func() {
		board.closeDoor()
	})
	return true
}

func (board *FakeBoard) closeDoor() bool {
	fmt.Println("FAKEBOARD: Close door")
	board.publisher.SendMessage("door.state.lock", "Door Locked")
	return true
}

func (board *FakeBoard) setPublisher(publisher Publisher) bool {
	fmt.Println("FAKEBOARD: Set publisher")
	board.publisher = publisher
	return true
}
