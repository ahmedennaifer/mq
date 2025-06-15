package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	From      string `json:"from"`
	Topic     string `json:"topic"`
}

func NewMessage(message string) *Message {
	return &Message{
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}
}

func Serialize(m Message) (error, []byte) {
	serialized, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("error serializing message: %v", m.Message)
		return err, nil
	}
	return nil, serialized
}

func Deserialize(m string) (error, *Message) {
	// m is of type string, because it will be convreted from buff to string in parseIntoCommand
	var msg *Message
	err := json.Unmarshal([]byte(m), &msg)
	if err != nil {
		fmt.Printf("error deserializing message: %v", string(m))
		return err, nil
	}
	return nil, msg
}
