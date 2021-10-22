package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// this is magic
// it is global - all clients shares it
var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error on connection upgradation:", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("error on reading message:", err)
			break
		}
		log.Printf("received: %s", msg)
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("error on writing message:", err)
			break
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func main() {
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
