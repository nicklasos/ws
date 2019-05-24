package main

import (
	"testing"
)

func TestParseMessage(t *testing.T) {
	msg, err := ParseChatMessage([]byte("[\"chat\", \"room\", \"message text\"]"))

	if err != nil {
		t.Fatal("Error on parse message")
	}

	if msg.messageType != "chat" {
		t.Error("Message type should be chat")
	}

	if msg.room != "room" {
		t.Error("Message room should be room")
	}

	if msg.message != "message text" {
		t.Error("Message text is wrong")
	}
}
