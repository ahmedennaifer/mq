package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Type         string
	Addr         string
	Port         string
	Ln           net.Listener
	Topics       []Topic
	WaitingPeers []Peer
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
	for {
		// accept loop
		conn, err := s.Ln.Accept()
		if err != nil {
			fmt.Printf("cannot accept connection from %v: %v", conn.RemoteAddr(), err)
			return err
		}
		peer := NewPeer(conn.RemoteAddr().String(), conn)
		s.WaitingPeers = append(s.WaitingPeers, *peer)

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
		match, err := s.GetTopic(cmd.Payload)

		if err != nil {
			fmt.Printf("error retrieving topic: %v\n", err)
			return err
		}
		// TODO: check if peer exists before asign

		waitingPeer, err := s.GetWaitingPeer(cmd.Target)
		if err != nil {
			fmt.Printf("coudlnt find waiting peer\n")
			return err
		}
		err = waitingPeer.Subscribe(cmd.Payload)
		if err != nil {
			fmt.Printf("error sub: %v\n", err)
		}
		if _, err := waitingPeer.Conn.Write([]byte("success")); err != nil {
			fmt.Printf("Error sending res to client %v", err)
		}

		match.Peers = append(match.Peers, *waitingPeer)
		fmt.Printf("added peer to topic\n")

		// peer no longer waiting, we remove.
		for i, peer := range s.WaitingPeers {
			if peer.Name == waitingPeer.Name {
				s.WaitingPeers = append(s.WaitingPeers[:i], s.WaitingPeers[i+1:]...)
				fmt.Print("removed from waiting\n")
				break
			}
		}

	case "waiting":
		peerMap := make(map[string]int)
		var peers []string

		for _, peer := range s.WaitingPeers {
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

func (s *Server) GetWaitingPeer(name string) (*Peer, error) {
	for _, peer := range s.WaitingPeers {
		if peer.Name == name {
			return &peer, nil
		}
	}
	return &Peer{}, errors.New("peer is not in waiting list")
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
	fmt.Printf("%v\n", name)
	for _, topic := range s.Topics {
		if topic.Name == name {
			return &topic, nil
		}
	}
	return &Topic{}, errors.New("topic does not exist\n")
}

func main() {
	conf := ServerConfig{
		addr: "localhost",
		port: "8080",
	}
	server := NewServer(conf)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
