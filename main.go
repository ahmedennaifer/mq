package main

import (
	"flag"
	"fmt"
	"log"
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
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
