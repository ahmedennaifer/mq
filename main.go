package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	mode := flag.String("mode", "server", "server or client")
	flag.Parse()

	switch *mode {
	case "server":
		startServer()
	case "client":
		startClient()
	default:
		fmt.Println("Usage: go run . -mode=server or go run . -mode=client")
	}
}

func startServer() {
	conf := ServerConfig{
		addr: "localhost",
		port: "8080",
	}
	server := NewServer(conf)
	topic := NewTopic("test-topic")
	err := server.AddTopic(topic.Name)
	if err != nil {
		fmt.Printf("error creating topic: %v", err)
	}

	fmt.Println("Starting server on localhost:8080...")
	go func() {
		err = server.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Second * 20)
	fmt.Print("sleeping for 20 sec...\n")
	for i := range 100 {
		fmt.Print("broadcasting...\n")
		msg := fmt.Sprintf("broadcasting test: %v", i+1)
		topic.Broadcast(msg)
	}
}
