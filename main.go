package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error on loading .env file")
	}

	hub := newHub()
	go hub.run()

	stackInit()
	queueInit()
	defer queueShutdown()

	go onlineInit(hub)
	go queueRun(hub)
	go runClearStats()

	serveHTTP(hub)
}
