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
	}
}

func startServer() {
	conf := ServerConfig{
		addr: "localhost",
		port: "8080",
	}
	server := NewServer(conf)
	err := server.AddTopic("test-topic")
	if err != nil {
		fmt.Printf("error creating topic: %v", err)
	}

	fmt.Println("Starting server on localhost:8080...")
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
