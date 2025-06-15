package main

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
)

type Peer struct {
	Name     string
	Conn     net.Conn
	Messages []string
	Topics   []string
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Name:     conn.RemoteAddr().String(),
		Conn:     conn,
		Messages: make([]string, 1),
	}
}

func (p *Peer) Subscribe(topic string) error {
	if slices.Contains(p.Topics, topic) {
		fmt.Println("Already subbed to topic", topic)
		return errors.New("error: already subscribed to topic")
	}

	sconn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("error connecting", err)
		return err
	}
	defer sconn.Close()

	subCmdStr := fmt.Sprintf("subscribe %v %v", p.Name, topic)
	_, err = sconn.Write([]byte(subCmdStr))
	if err != nil {
		fmt.Println("error sending subcmdstr to server", err)
		return err
	}

	readBuff := make([]byte, 1024)
	n, err := sconn.Read(readBuff)
	if err != nil {
		fmt.Println("error reading response for str cmd parsing", err)
		return err
	}

	response := strings.TrimSpace(string(readBuff[:n]))
	if response != "success" {
		fmt.Printf("failed to subscribe, got: %v\n", response)
		return errors.New("failed to subscribe")
	}

	p.Topics = append(p.Topics, topic)
	fmt.Print("subbed with success\n")
	fmt.Printf("peer topics: %v\n", p.Topics)
	return nil
}
