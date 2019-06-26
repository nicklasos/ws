package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/tomasen/realip"

	"github.com/gorilla/websocket"
	"github.com/nicklasos/golimiter"
)

var (
	ipLimit = golimiter.NewLimiter(5, 5)
	idLimit = golimiter.NewLimiter(2, 3)

	banIpTime = 30 * time.Minute
	banIdTime = 20 * time.Minute
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Limit rooms that users can join per connection
	maxRooms = 2
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  1024,
	// WriteBufferSize: 1024,
	// @todo: fix
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Join time
	time int64

	// User id
	id string

	// User ip
	ip string

	// Join rooms
	rooms []string

	// Buffered channel of outbound messages.
	send chan []byte
}

// InitParams is params from client (browser, etc)
type InitParams struct {
	id    string
	rooms []string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// log.Println("Message from client: ", string(message))

		chatMsg, err := parseChatMessage(c, message)
		if err != nil {
			log.Println("Error on parsing chat message from client", err)
			break
		}

		if isSpam(chatMsg.message) {
			ipLimit.Ban(c.ip, banIpTime)
			c.hub.unregister <- c
			break
		}

		if !idLimit.Allow(c.id) {
			idLimit.Ban(c.id, banIdTime)
			c.hub.unregister <- c
			break
		}
		if !ipLimit.Allow(c.ip) {
			ipLimit.Ban(c.ip, banIpTime)
			c.hub.unregister <- c
			break
		}

		logChatMessage(chatMsg)

		c.hub.chat <- chatMsg
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func parseInitParams(values map[string][]string) (*InitParams, error) {
	id, ok := values["id"]
	if !ok || len(id[0]) < 1 {
		return nil, errors.New("Param id is missing")
	}

	rooms, _ := values["rooms"]
	if len(rooms) > maxRooms {
		return nil, errors.New("Param rooms is out of limit")
	}

	return &InitParams{id[0], rooms}, nil
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	initParams, err := parseInitParams(r.URL.Query())
	if err != nil {
		log.Println(err)
		return
	}

	ip := realip.FromRequest(r)

	if ipLimit.IsBanned(ip) || idLimit.IsBanned(initParams.id) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("429 - Too many requests"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// log.Println("Initial params", initParams)

	client := &Client{
		hub:   hub,
		conn:  conn,
		time:  time.Now().Unix(),
		id:    initParams.id,
		ip:    realip.FromRequest(r),
		rooms: initParams.rooms,
		send:  make(chan []byte, 256),
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
