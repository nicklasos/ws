package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func OnlineInit(hub *Hub) {
	message := make([]interface{}, 2)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:

				message[0] = "online"
				users := GetStats(hub).Users
				message[1] = users

				StackPutKeyValues("ws.online", users)

				msg, err := json.Marshal(message)
				if err != nil {
					fmt.Println("Error on json encode online stats", err)
				} else {
					hub.broadcast <- msg
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
