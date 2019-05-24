package main

import (
	"encoding/json"
	"errors"
)

// ChatMessage is parsed json from client
type ChatMessage struct {
	messageType string
	room        string
	message     string
	messageRaw  []byte
	client      *Client
}

// ParseChatMessage parses raw json message from client
func ParseChatMessage(message []byte) (*ChatMessage, error) {
	var data []string

	err := json.Unmarshal(message, &data)
	if err != nil {
		return nil, err
	}

	if len(data) != 3 {
		return nil, errors.New("Wrong message from client")
	}

	if data[0] != "chat" {
		return nil, errors.New("Wrong message type for chat")
	}

	return &ChatMessage{data[0], data[1], data[2], message, nil}, nil
}
