package main

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
