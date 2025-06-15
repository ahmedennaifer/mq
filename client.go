package main

import (
	"fmt"
	"net"
	"time"
)

func startClient() {
	fmt.Println("Starting client...")
	time.Sleep(500 * time.Millisecond)
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	peer := &Peer{
		Name:     conn.LocalAddr().String(),
		Conn:     conn,
		Messages: make([]string, 1),
		Topics:   make([]string, 0),
	}

	err = peer.Subscribe("test-topic")
	if err != nil {
		fmt.Printf("Subscribe error: %v\n", err)
	} else {
		fmt.Printf("Peer %v subscribed to topics: %v\n", peer.Name, peer.Topics)
	}
	time.Sleep(2 * time.Second)
}
