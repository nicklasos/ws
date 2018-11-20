package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error on loading .env file")
	}

	hub := newHub()
	go hub.run()

	StackInit()
	OnlineInit(hub)
	QueueInit()
	defer QueueShutdown()

	go QueueRun(hub)

	ServeHTTP(hub)
}
