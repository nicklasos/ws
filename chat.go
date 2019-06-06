package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
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
	incrementKey(fmt.Sprintf("%s.%s", "messages.count", msg.room))
	setAdd(fmt.Sprintf("%s.%s", "messages.writers", msg.room), msg.client.id)
}

type ChatRoomStats struct {
	messages int64
	writets  int64
}

func runClearStats() {
	for {
		time.Sleep(time.Hour * 24 * 10)
		logClearStats()
	}
}

func logClearStats() {
	keys := getKeys("messages.writers.*")
	countKeys := getKeys("messages.count.*")

	for _, countKey := range countKeys {
		keys = append(keys, countKey)
	}

	for _, key := range keys {
		keyDel(key)
	}
}

func logGetStats(room string) *ChatRoomStats {
	writets := setCount(fmt.Sprintf("%s.%s", "messages.writers", room))
	msgStr := getKey(fmt.Sprintf("%s.%s", "messages.count", room))

	msgCnt, err := strconv.ParseInt(msgStr, 10, 64)
	if err != nil {
		msgCnt = 0
	}

	return &ChatRoomStats{msgCnt, writets}
}
