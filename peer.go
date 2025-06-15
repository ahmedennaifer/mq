package main

import (
	"errors"
	"fmt"
	"net"
	"slices"
)

type Peer struct {
	Name     string
	Conn     net.Conn
	Messages []string
	Topics   []string
}

func NewPeer(name string, conn net.Conn) *Peer {
	return &Peer{
		Name:     name,
		Conn:     conn,
		Messages: make([]string, 1),
	}
}

func (p *Peer) Subscribe(topic string) error {
	// Peer sends sub request.
	// server parses request
	// server adds peer to topics
	// server returns OK
	// try first with cli, then switch to lib. TODO
	if slices.Contains(p.Topics, topic) {
		fmt.Println("Already subbed to to topic", topic)
		return errors.New("error: topic already exists")
	}
	sconn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("error connecting", err)
	}
	subCmdStr := fmt.Sprintf("subscribe %v %v", p.Name, topic)
	_, err = sconn.Write([]byte(subCmdStr))
	if err != nil {
		fmt.Println("error sending subcmdstr to server", err)
		return err
	}
	var readBuff []byte
	n, err := p.Conn.Read(readBuff)
	if err != nil {
		fmt.Println("error reading response for str cmd parsing", err)
		return err
	}
	if string(readBuff[:n]) != "success" {
		fmt.Printf("failed to subscribe, got: %v", string(readBuff[:n]))
		return errors.New("failed to subscribe")
	}
	p.Topics = append(p.Topics, topic)
	fmt.Print("subbed with success\n")
	fmt.Printf("peer topics : %v", p.Topics)
	return nil

}
