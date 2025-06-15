package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"slices"
)

type Server struct {
	Type   string
	Addr   string
	Port   string
	Ln     net.Listener
	Topics []Topic
	Peers  map[net.Conn]*Peer
}

type ServerConfig struct{ addr, port string }

func NewServer(config ServerConfig) *Server {
	return &Server{
		Type: "tcp",
		Addr: config.addr,
		Port: config.port,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen(s.Type, s.Addr+":"+s.Port)
	if err != nil {
		log.Fatalf("error starting server %v\n", err)
		return err
	}
	s.Ln = ln
	s.Peers = make(map[net.Conn]*Peer)

	for {
		// accept loop
		conn, err := s.Ln.Accept()
		if err != nil {
			fmt.Printf("cannot accept connection from %v: %v", conn.RemoteAddr().String(), err)
			return err
		}
		peer := NewPeer(conn)
		s.Peers[conn] = peer

		go s.handleConnection(conn)

	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// peer handled here
	fmt.Printf("client %v connected\n", conn.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		// read loop
		n, err := conn.Read(buf)
		// replace with protocol parsing into peer+action
		if err != nil {
			fmt.Println("error reading msg\n", err)
			return
		}

		cmd, err := parseIntoCommand(buf[:n])
		if err != nil {
			fmt.Printf("Couldnt parse response: %v\n", err)
			if _, err := conn.Write([]byte("\033[31m" + err.Error() + "\033[0m\n")); err != nil {
				fmt.Println("cannot write to client", err)
			}
		}
		if cmd.Action == "" {
			fmt.Printf("%v: %v\n", conn.RemoteAddr(), string(buf[:n]))
		} else {
			if handleErr := s.handleCommand(conn, *cmd); err != nil {
				fmt.Println("error:", handleErr)
				if _, err := conn.Write([]byte(handleErr.Error())); err != nil {
					fmt.Println("cannot write to client", err)
				}
			}
			fmt.Printf("cmd: %v %v %v\n", cmd.Action, cmd.Target, cmd.Payload)
		}

	}
}

func (s *Server) handleCommand(conn net.Conn, cmd Command) error {
	switch cmd.Action {
	case "create":
		if cmd.Target == "" || cmd.Payload == "" {
			fmt.Println("target or payload unspecified!")
			return errors.New("error: target and/or payload must not be empty")
		}
		if err := s.AddTopic(cmd.Payload); err != nil {
			fmt.Println("error:", err)
			return err
		}

		if _, err := conn.Write([]byte("\033[32mTopic added with success\033[0m\n")); err != nil {
			fmt.Println("err", err)
			return err
		}
	case "list":

		var topics []string
		for _, topic := range s.Topics {
			topics = append(topics, topic.Name)
		}
		str := "\033[32m" + fmt.Sprintf("%v", topics) + "\033[32m"

		if _, err := conn.Write([]byte(str)); err != nil {
			fmt.Printf("error: couldnt send list to client: %v\n", err)
			return err
		}

	case "broadcast":
		match, err := s.GetTopic(cmd.Target)
		if err != nil {
			fmt.Printf("error retrieving topic: %v\n", err)
			return err
		}
		if err := match.Broadcast("hello this is a broadcast\n"); err != nil {
			fmt.Printf("Error broadcasting message to clients %v\n", err)
			return err
		}

	case "subscribe":
		topic, err := s.GetTopic(cmd.Payload)
		if err != nil {
			conn.Write([]byte("error: topic not found"))
			return err
		}

		peer, exists := s.Peers[conn]
		if !exists {
			conn.Write([]byte("error: peer not registered"))
			return errors.New("peer not found")
		}

		if slices.Contains(peer.Topics, cmd.Payload) {
			conn.Write([]byte("error: already subscribed"))
			return errors.New("already subscribed")
		}

		peer.Topics = append(peer.Topics, cmd.Payload)
		topic.Peers = append(topic.Peers, *peer)

		conn.Write([]byte("success"))
		fmt.Printf("Peer %s joined topic %s\n", peer.Name, cmd.Payload)
		return nil

	case "peers":
		peerMap := make(map[string]int)
		var peers []string

		for _, peer := range s.Peers {
			topicCount := len(peer.Topics)
			peerMap[peer.Name] = topicCount
			peerInfo := fmt.Sprintf("%s (%d topics)", peer.Name, topicCount)
			peers = append(peers, peerInfo)
		}

		str := "\033[32m" + fmt.Sprintf("%v", peers) + "\033[0m"
		if _, err := conn.Write([]byte(str)); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (s *Server) AddTopic(name string) error {
	for _, topic := range s.Topics {
		if topic.Name == name {
			fmt.Printf("error: topic %v already exists!\n", name)
			return errors.New("error: topic already exists")
		}
	}
	t := NewTopic(name)
	s.Topics = append(s.Topics, *t)
	return nil
}

func (s *Server) GetTopic(name string) (*Topic, error) {
	for i, topic := range s.Topics {
		if topic.Name == name {
			return &s.Topics[i], nil
		}
	}
	return &Topic{}, errors.New("topic does not exist")
}
