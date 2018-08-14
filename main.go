package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"encoding/json"
	"github.com/joho/godotenv"
)

var tmpl = template.Must(template.ParseFiles("websockets.html"))

type TmplData struct {
	Port   string
	Schema string
}

type Stats struct {
	Connections int `json:"connections"`
	Users       int `json:"users"`
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error on loading .env file")
	}

	log.Println("DEBUG", os.Getenv("DEBUG"))

	hub := newHub()
	go hub.run()

	StackInit()
	QueueInit()
	defer QueueShutdown()

	go QueueRun(hub)

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
		uniq := make(map[string]bool)
		for client := range hub.clients {
			uniq[client.id] = true
		}

		stats := Stats{len(hub.clients), len(uniq)}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(stats)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if os.Getenv("USE_SSL") == "true" {
		log.Println("Start server https")

		if err := http.ListenAndServeTLS(os.Getenv("DOMAIN")+":"+os.Getenv("PORT"), "server.crt", "server.key", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		log.Println("Start server http")

		if err := http.ListenAndServe(os.Getenv("DOMAIN")+":"+os.Getenv("PORT"), nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}

}
