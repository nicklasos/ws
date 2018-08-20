package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl = template.Must(template.ParseFiles("websockets.html"))

type TmplData struct {
	Port   string
	Schema string
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tmplData := TmplData{
		os.Getenv("TMPL_PORT"),
		"ws",
	}

	if os.Getenv("USE_SSL") == "true" {
		tmplData.Schema = "wss"
	}

	tmpl.Execute(w, tmplData)
}

func send(hub *Hub, w http.ResponseWriter, r *http.Request) {
	hub.broadcast <- []byte("from http")

	w.Write([]byte("Message sent"))
}

func ServeHTTP(hub *Hub) {
	if os.Getenv("DEBUG") == "true" {
		http.HandleFunc("/", serveHome)

		http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
			send(hub, w, r)
		})

		http.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
			QueueSend()
			w.Write([]byte("Message sent to queue"))
		})

		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	}

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		Stats(hub, w)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	address := os.Getenv("DOMAIN") + ":" + os.Getenv("PORT")

	if os.Getenv("USE_SSL") == "true" {
		log.Println("Start server https")

		if err := http.ListenAndServeTLS(address, "server.crt", "server.key", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		log.Println("Start server http")

		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
