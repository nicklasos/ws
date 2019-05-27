package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ChatMessage is parsed json from client
type ChatMessage struct {
	messageType string
	room        string
	message     string
	messageRaw  []byte
	client      *Client
}

// parseChatMessage parses raw json message from client
func parseChatMessage(client *Client, message []byte) (*ChatMessage, error) {
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

	return &ChatMessage{data[0], data[1], data[2], message, client}, nil
}

func logChatMessage(msg *ChatMessage) {
	stackPush(fmt.Sprintf("%s.%s", msg.messageType, msg.room), 30, msg.messageRaw)
}
