package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var tmpl = template.Must(template.ParseFiles("websockets.html"))

type TmplData struct {
	Port   string
	Schema string
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
	}

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
