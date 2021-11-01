package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	wsmsg "websocket-client/wsmsg"
)

var (
	done      chan interface{}
	interrupt chan os.Signal
	sendMsg   chan []byte
)

var (
	sessionIDP *string
	usernameP  *string
)

const (
	TypeConnected = "connected"
	TypeJoin      = "join-session"
	TypeMsg       = "message"
)

func receiveHandler(conn *websocket.Conn) {
	defer close(done)
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Println("error on receiving:", err)
			return
		}

		msg := &wsmsg.WebsocketMessage{}
		if err := proto.Unmarshal(b, msg); err != nil {
			log.Printf("err: %v. %v: %v", err, msg.Username, msg.Message)
		}

		// filtering and displaying
		switch msg.Type {
		case TypeConnected:
			log.Printf("connected: %v", msg.Message)
			// ask for join automatically
			msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
				Type:      TypeJoin, // TODO: proto enum
				SessionId: *sessionIDP,
				Username:  *usernameP,
				// Message:   fmt.Sprintf("Let me join the session named: %v", *sessionIDP),
			})
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}

			sendMsg <- msgBytes
		case TypeMsg:
		default:
			log.Printf("unknown type of message: %+v", msg.Type)
		}
		fmt.Printf("msg received: %v: %v\n", msg.Username, msg.Message)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlP := flag.String("url", "", "websocket server url")
	sessionIDP = flag.String("sid", "", "roodid to enter")
	usernameP = flag.String("uname", "", "username to use")
	flag.Parse()

	done = make(chan interface{})
	interrupt = make(chan os.Signal)
	sendMsg = make(chan []byte)

	signal.Notify(interrupt, os.Interrupt)

	socketURL := "ws://localhost:8080" + "/ws"
	if *urlP != "" {
		socketURL = *urlP
	}
	log.Printf("websocket server URL: %v", socketURL)

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

			msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
				Type:      TypeMsg, // TODO: proto enum
				SessionId: *sessionIDP,
				Username:  *usernameP,
				Message:   message,
			})
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}

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
