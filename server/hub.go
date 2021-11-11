// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	wsmsg "plantumlive-ws-server/wsmsg"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*websocket.Conn]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients. - just connect
	register chan *Client

	// Join messages from the clients - after connect
	join chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 10),
		register:   make(chan *Client),
		join:       make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*websocket.Conn]*Client),
	}
}

const (
	TypeConnected = "connected"
	TypeJoin      = "joinSession"
	TypeMsg       = "message"
)

const messageDir = "./messages/"
const serverAnnounce = "server"

func saveSessionMessage(msg *wsmsg.WebsocketMessage, b []byte) error {

	// marshalled protobuf message print
	log.Printf("save messageBytes: %x", b)

	// open file
	fn := messageDir + msg.SessionId + ".message"
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// composite saving format
	// type(1 byte) + lenght(4 bytes) + b
	msgType := []byte{byte(msg.Type)}

	length := make([]byte, 4)
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, int32(len(b))); err != nil {
		return err
	}
	copy(length, buf.Bytes())
	b = append(length, b...)
	b = append(msgType, b...)

	log.Printf("formed messageBytes: %x", b)

	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

// func getSessionMessage(msg wsmsg.WebsocketMessage) (messages []wsmsg.WebsocketMessage) {

// 	fn := messageDir + msg.SessionId + ".message"

// 	file, err := os.Open(fn)
// 	if err != nil {
// 		fmt.Println(err)
// 		return messages
// 	}
// 	defer file.Close()

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(file)

// 	bytes := buf.Bytes()
// 	for len(bytes) > 5 {
// 		// read type - do nothing for now
// 		// msgType := bytes[0]
// 		length := binary.LittleEndian.Uint32(bytes[1:5])
// 		msgBytes := bytes[5 : 5+length]
// 		bytes = bytes[6+length:]

// 		// convert to wsmsg.WebsocketMessage and append
// 		msg := wsmsg.WebsocketMessage{}
// 		if err := proto.Unmarshal(msgBytes, &msg); err != nil {
// 			log.Println("fail to unmarshal:", err)
// 			continue
// 		}
// 		messages = append(messages, msg)
// 	}
// 	return messages
// }
func getSessionMessageBytes(msg *wsmsg.WebsocketMessage) (msgBytesSlice [][]byte) {

	fn := messageDir + msg.SessionId + ".message"

	file, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return msgBytesSlice
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	bytes := buf.Bytes()
	log.Printf("read file len: %v", len(bytes))
	for {
		// read type - do nothing for now
		// msgType := bytes[0]
		length := binary.LittleEndian.Uint32(bytes[1:5])
		log.Println("length:", length)

		msgBytes := bytes[5 : 5+length]
		restBytesLen := len(bytes) - len(msgBytes)
		if restBytesLen < 5 {
			log.Printf("left bytes: %x", restBytesLen)
			break
		}
		bytes = bytes[5+length:]

		msgBytesSlice = append(msgBytesSlice, msgBytes)
	}
	return msgBytesSlice
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Println("an user connected to websocket server")
			h.clients[client.conn] = client

			welcomeMsg := &wsmsg.WebsocketMessage{
				Type:    wsmsg.Type_Connected,
				Message: "You are connected to websocket server",
			}

			msgBytes, err := proto.Marshal(welcomeMsg)
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}
			client.send <- msgBytes

		case client := <-h.unregister:
			if leaveClient, ok := h.clients[client.conn]; ok {
				fmt.Printf("disconnect user %s in the session %s\n", client.Username, client.SessionID)

				// leave message to everybody
				msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
					Type:      wsmsg.Type_Join,
					SessionId: leaveClient.SessionID,
					Username:  leaveClient.Username,
					Message:   fmt.Sprintf("=== %s just leaved ===", leaveClient.Username),
				})
				if err != nil {
					log.Println("fail to marshal:", err)
					continue
				}
				h.broadcast <- msgBytes

				// delete user
				delete(h.clients, client.conn)
				close(client.send)
			}

		case joinClient := <-h.join:
			// re-register client with SessionId and Username
			h.clients[joinClient.conn] = joinClient

			// join message to everybody
			msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
				Type:      wsmsg.Type_Join,
				SessionId: joinClient.SessionID,
				Username:  joinClient.Username,
				Message:   fmt.Sprintf("=== %s just joined ===", joinClient.Username),
			})
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}
			h.broadcast <- msgBytes

		case msgBytes := <-h.broadcast:

			// filtering: check the SessionId and Username
			msg := &wsmsg.WebsocketMessage{}
			if err := proto.Unmarshal(msgBytes, msg); err != nil {
				log.Println("fail to unmarshal:", err)
				continue
			}
			fmt.Printf("session %s: user %s\t: %v\n", msg.SessionId, msg.Username, msg.Message)
			if err := saveSessionMessage(msg, msgBytes); err != nil {
				log.Println("fail to saveSassionMessage:", err)
			}

			for _, client := range h.clients {
				if client.Username == msg.Username {
					if msg.Type == wsmsg.Type_Join {
						log.Printf("send all history first")
						// before sending welcome message, send all the history to newlyjoined
						// getSessionMessage(msg) []wsmsg.WebsocketMessage and loop sending
						messageHistory := getSessionMessageBytes(msg)
						log.Printf("message count: %v", len(messageHistory))
						for _, messageBytes := range messageHistory {
							log.Printf("messageBytes: %x", messageBytes)
							select {
							case client.send <- messageBytes:
							default: // if client cannot consume, clear client
								close(client.send)
								delete(h.clients, client.conn)
							}
						}
						continue
					}
				}

				if client.SessionID != msg.SessionId {
					log.Printf("send msg only to proper session clients: %s != %s", client.SessionID, msg.SessionId)
					continue
				}

				select {
				case client.send <- msgBytes:
				default: // if client cannot consume, clear client
					close(client.send)
					delete(h.clients, client.conn)
				}
			}
		}
	}
}
