package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
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

func receiveHandler(conn *websocket.Conn) {
	defer close(done)
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Println("error on receiving:", err)
			return
		}
		log.Printf("received %x", b)
		msg := &wsmsg.WebsocketMessage{}
		if err := proto.Unmarshal(b, msg); err != nil {
			log.Printf("err: %v. %v: %v", err, msg.Username, msg.Message)
		}

		// filtering and displaying
		switch msg.GetType() {
		case wsmsg.Type_Connected:
			fmt.Printf("=== websocket server connected: %v ===\n", msg.Message)
			// ask for join automatically
			msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
				Type:      msg.GetType(), // TODO: proto enum
				SessionId: *sessionIDP,
				Username:  *usernameP,
				Message:   fmt.Sprintf("Let me join the session named: %v", *sessionIDP),
			})
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}
			sendMsg <- msgBytes
		case wsmsg.Type_Join:
			if msg.GetUsername() == *usernameP {
				fmt.Printf("=== joined to the sesson %s successfully ===\n", msg.GetSessionId())
			} else {
				fmt.Printf("\n%v\n\n", msg.GetMessage())
			}
		case wsmsg.Type_Msg:
			fmt.Printf("session %s: user %s\t: %v\n", msg.SessionId, msg.Username, msg.Message)
		default:
			fmt.Printf("unknown type of message: %s\n", msg.GetType())
		}
		log.Printf("msg received: %v: %v\n", msg.Username, msg.Message)
	}
}

var (
	urlP       = flag.String("url", "", "websocket server url. default is ws://localhost:8080/ws")
	sessionIDP = flag.String("s", "", "sessionid to join")
	usernameP  = flag.String("u", "", "username to use")
	debug      = flag.Bool("d", false, "display debugging log")
)

func scanMessage() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("session id: %v, username: %v", *sessionIDP, *usernameP)
		msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
			Type:      wsmsg.Type_Msg,
			SessionId: *sessionIDP,
			Username:  *usernameP,
			Message:   message,
		})
		if err != nil {
			log.Println("fail to marshal:", err)
			continue
		}
		// before send message, rewrite the scanned msg
		sendMsg <- msgBytes
	}
}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !(*debug) {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("session id: %v, username: %v", *sessionIDP, *usernameP)

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
	go scanMessage()

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
