package main

import (
	"log"
	"net"
)

var (
	EventChannel = make(chan uint8)
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Print(err)
		return
	}
	RunLoop(ln)
}

func RunLoop(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go Dispatch(conn, EventChannel)
	}
}

const (
	Open  = 1
	Close = 4
)

func Dispatch(conn net.Conn, events chan<- uint8) {
	defer conn.Close()
	events <- Close
}
