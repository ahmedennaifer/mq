package main

import "net"

type Peer struct {
	Name     string
	conn     net.Conn
	Messages []string
}

func NewPeer(name string, conn net.Conn) *Peer {
	return &Peer{
		Name:     name,
		conn:     conn,
		Messages: make([]string, 1),
	}
}
