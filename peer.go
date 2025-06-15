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

func (p *Peer) Publish(topic string, msg *Message) error { // TODO: message struct {
	// p.Publish(topic, message)
	// write message to topic messages
	err, serialize := Serialize(*msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	msg.From = p.Name
	msg.Topic = topic

	// | for avoiding splitting on json parts. ie: user:< >name,< >timestamp: and so on.
	cmd := fmt.Sprintf("publish %s|%s", serialize, msg.Topic)

	_, err = p.Conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send join command: %v", err)
	}
	fmt.Printf("Publishing to %v: %v", msg.Message, topic)

	buf := make([]byte, 1024)
	_, err = p.Conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	return nil
}

func (p *Peer) Subscribe(topic string) error {
	if slices.Contains(p.Topics, topic) {
		return errors.New("already subscribed to topic")
	}

	cmd := fmt.Sprintf("subscribe %v %v", p.Name, topic)
	_, err := p.Conn.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send join command: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := p.Conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	response := strings.TrimSpace(string(buf[:n]))
	if response != "success" {
		return fmt.Errorf("subscription failed: %s", response)
	}

	p.Topics = append(p.Topics, topic)
	return nil
}
