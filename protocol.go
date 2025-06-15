package main

import (
	"errors"
	"fmt"
	"strings"
)

// const (
// CommandHealth string
// CommandNewTopic string
// CommandShowTopics string
// CommandPublish string
// )

type Command struct {
	Action  string // create, health, sub, unsub, ..
	Target  string // action target name
	Payload string // msg
	// From Peer
	// To Peer
}

func parseIntoCommand(buff []byte) (*Command, error) {
	if len(buff) == 0 {
		fmt.Println("error: buffer is empty! ")
		return &Command{}, errors.New("empty buff")
	}
	var (
		stripped = strings.TrimSpace(string(buff))
		res      = strings.Split(stripped, " ")
	)

	switch res[0] {

	case "topic": // 5 bytes
		if len(res) < 2 {
			fmt.Printf("error: arguments must be exactly 2. got: %v\n", len(res))
			return &Command{}, errors.New("must provide topic name")
		}
		return &Command{
			Action:  "create",
			Target:  "topic",
			Payload: res[1],
		}, nil

	case "list": // 5 bytes
		return &Command{
			Action:  "list",
			Target:  "topic",
			Payload: "",
		}, nil

	case "broadcast":
		if len(res) < 3 {
			fmt.Printf("error: arguments must be exactly 3. got: %v\n", len(res))
			return &Command{}, errors.New("must provide topic name and message")
		}
		return &Command{
			Action:  "broadcast",
			Target:  res[1],
			Payload: res[2],
		}, nil

	case "subscribe":
		if len(res) < 3 {
			fmt.Printf("error: arguments must be exactly 3. got: %v\n", len(res))
			return &Command{}, errors.New("must provide topic name and message")
		}
		return &Command{
			Action:  "subscribe",
			Target:  res[1],
			Payload: res[2],
		}, nil

	case "peers":
		return &Command{
			Action:  "peers",
			Target:  "",
			Payload: "",
		}, nil

	case "publish":
		// split on | for publish: "publish <json>|<topic>"
		parts := strings.SplitN(stripped[8:], "|", 2)
		if len(parts) != 2 {
			return &Command{}, errors.New("publish format: publish <json>|<topic>")
		}

		return &Command{
			Action:  "publish",
			Target:  parts[1], // topic
			Payload: parts[0], // json
		}, nil

	default:
		return &Command{}, nil
	}
}
