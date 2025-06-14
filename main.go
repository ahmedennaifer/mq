package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

// 1 - tcp accept/read loop
// 2 - message/peer struct
// 3 - broadcast

type Server struct {
	Type   string
	Addr   string
	Port   string
	Ln     net.Listener
	Topics []Topic
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
		}
		if cmd.Action == "" {
			fmt.Printf("%v: %v\n", conn.RemoteAddr(), string(buf[:n]))
		} else {
			if err := s.handleCommand(conn, *cmd); err != nil {
				fmt.Println("error:", err)
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
			return errors.New("error: couldnt add topic")
		}

		if _, err := conn.Write([]byte("\033[32mTopic added with success\033[0m\n")); err != nil {
			fmt.Println("err", err)
			return errors.New("error: couldnt send response to client")
		}
	case "list":

		var topics []string
		for _, topic := range s.Topics {
			topics = append(topics, topic.Name)
		}
		str := "\033[32m" + fmt.Sprintf("%v", topics) + "\033[32m"

		if _, err := conn.Write([]byte(str)); err != nil {
			fmt.Printf("error: couldnt send list to client: %v\n", err)
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
