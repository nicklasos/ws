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

func serveHTTP(hub *Hub) {
	if os.Getenv("DEBUG") == "true" {
		http.HandleFunc("/", serveHome)

		http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
			send(hub, w, r)
		})

		http.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
			queueSend()
			w.Write([]byte("Message sent to queue"))
		})

		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
	}

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats(hub, w)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	address := os.Getenv("DOMAIN") + ":" + os.Getenv("PORT")

	if os.Getenv("USE_SSL") == "true" {
		log.Println("Start https server on port: ", address)

		if err := http.ListenAndServeTLS(address, "server.crt", "server.key", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		log.Println("Start http server on port: ", address)

		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
