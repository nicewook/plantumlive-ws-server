package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
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

var (
	clientRoomID   string
	clientUsername string
)

type WebsocketMessage struct {
	RoomID   string `json:"roomid,omitempty"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message,omitempty"`
}

func receiveHandler(conn *websocket.Conn) {
	defer close(done)
	var msg WebsocketMessage
	for {
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("error on receiving:", err)
			return
		}

		if msg.Message == "welcome" {
			clientRoomID = msg.RoomID
			clientUsername = msg.Username
			fmt.Printf("Welcome %v", msg.Username)
			continue
		}
		log.Printf("%v: %v", msg.Username, msg.Message)
	}
}

func main() {

	roomIDPtr := flag.String("roomid", "", "rood id to enter")
	flag.Parse()

	done = make(chan interface{})
	interrupt = make(chan os.Signal)
	sendMsg := make(chan []byte)

	signal.Notify(interrupt, os.Interrupt)

	socketURL := "ws://localhost:8080" + "/ws"
	if *roomIDPtr != "" {
		socketURL = socketURL + "/" + *roomIDPtr
	}
	log.Println("socket url: ", socketURL)

	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	if err != nil {
		log.Fatal("error connecting to websocket server:", err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// client writes
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			message := scanner.Text()

			msgBytes, err := json.Marshal(WebsocketMessage{
				RoomID:   clientRoomID,
				Username: clientUsername,
				Message:  message,
			})
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}

			// RoomID, Username, Message - struct
			sendMsg <- msgBytes
		}
	}()

	for {
		select {
		case message := <-sendMsg:
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
