package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	done      chan interface{}
	interrupt chan os.Signal
)

func receiveHandler(conn *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("error on receiving:", err)
			return
		}
		log.Printf("received msg: %s", msg)
	}
}

func main() {
	done = make(chan interface{})
	interrupt = make(chan os.Signal)

	signal.Notify(interrupt, os.Interrupt)

	socketURL := "websocket://localhost:8080" + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	if err != nil {
		log.Fatal("error connecting to websocket server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	for {
		select {
		case <-time.After(1 * time.Second):
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from Client")); err != nil {
				log.Println("error on writing to websocket server:", err)
				return
			}
		case <-interrupt: // SIGINT(Ctrl+c)
			// graceful termination
			log.Println("SIGINT received. close all the connections")

			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("error on closing websocket:", err)
				return
			}

			select {
			case <-done:
				log.Println("receiver channel closed. exiting...")
			case <-time.After(1 * time.Second):
				log.Println("not waiting for receiving channel closed. exiting...")
			}
			return
		}
	}
}
