package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

type Rooms = map[string]*RoomStats
type ListOfRooms = []*RoomStats

type Data struct {
	Connections int         `json:"connections"`
	Users       int         `json:"users"`
	Users1min   int         `json:"users_1min"`
	Users5min   int         `json:"users_5min"`
	Users15min  int         `json:"users_15min"`
	Rooms       ListOfRooms `json:"rooms"`
	Version     string      `json:"version"`
	BannedUsers int         `json:"banned_users"`
	UptimeFrom  string      `json:"uptime_from"`
}

type RoomStats struct {
	Name     string `json:"name"`
	Online   int64  `json:"online"`
	Writers  int64  `json:"writers"`
	Messages int64  `json:"messages"`
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

	rooms := make(Rooms)

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
				rooms[room].Online++
			} else {
				r := logGetStats(room)
				rooms[room] = &RoomStats{room, 1, r.writets, r.messages}
			}
		}
	}

	listOfRooms := []*RoomStats{}
	for _, r := range rooms {
		listOfRooms = append(listOfRooms, r)
	}

	sort.Slice(listOfRooms, func(i, j int) bool {
		return listOfRooms[i].Online > listOfRooms[j].Online
	})

	return &Data{
		len(hub.clients),
		len(uniq),
		min1,
		min5,
		min15,
		listOfRooms,
		version,
		ipLimit.BannedCount() + idLimit.BannedCount(),
		uptime.Format("2006-01-02 15:04:05"),
	}
}

func stats(hub *Hub, w http.ResponseWriter) {
	s := getStats(hub)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}
