package main

import (
	"errors"
	"fmt"
)

type Topic struct {
	Name     string
	Peers    []Peer
	Messages []string
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name:     name,
		Peers:    make([]Peer, 0),
		Messages: make([]string, 0),
	}
}

func (t *Topic) Broadcast(payload string) error {
	fmt.Printf("Broadcasting to %v peers..\n", len(t.Peers))
	for _, peer := range t.Peers {
		if _, err := peer.Conn.Write([]byte(payload)); err != nil {
			fmt.Println(err)
			return errors.New("failed broadcasting to clients")
		}
	}
	return nil
}
