package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Data struct {
	Connections int            `json:"connections"`
	Users       int            `json:"users"`
	Users1min   int            `json:"users_1min"`
	Users5min   int            `json:"users_5min"`
	Users15min  int            `json:"users_15min"`
	Rooms       map[string]int `json:"rooms"`
}

func getStats(hub *Hub) *Data {
	now := time.Now().Unix()

	min1time := now - 60
	min5time := now - 60*5
	min15time := now - 60*15

	min1 := 0
	min5 := 0
	min15 := 0

	uniq := make(map[string]*Client)
	for client := range hub.clients {
		uniq[client.id] = client
	}

	rooms := make(map[string]int)

	for _, c := range uniq {
		if c.time > min1time {
			min1++
		}

		if c.time > min5time {
			min5++
		}

		if c.time > min15time {
			min15++
		}

		for _, room := range c.rooms {
			if _, ok := rooms[room]; ok {
				rooms[room]++
			} else {
				rooms[room] = 1
			}
		}
	}

	return &Data{
		Connections: len(hub.clients),
		Users:       len(uniq),
		Users1min:   min1,
		Users5min:   min5,
		Users15min:  min15,
		Rooms:       rooms,
	}
}

func stats(hub *Hub, w http.ResponseWriter) {
	s := getStats(hub)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}
