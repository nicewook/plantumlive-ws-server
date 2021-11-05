// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"log"
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
			if _, ok := h.clients[client.conn]; ok {
				fmt.Printf("disconnect user %s in the session %s\n", client.Username, client.SessionID)
				delete(h.clients, client.conn)
				close(client.send)
			}

		case joinClient := <-h.join:
			// re-register client with SessionId and Username
			h.clients[joinClient.conn] = joinClient

			// join message to everybody
			msgBytes, err := proto.Marshal(&wsmsg.WebsocketMessage{
				Type:      wsmsg.Type_Join, // TODO: proto enum
				SessionId: joinClient.SessionID,
				Username:  joinClient.Username,
				Message:   fmt.Sprintf("=== user %s just joined ===", joinClient.Username),
			})
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}
			h.broadcast <- msgBytes

		case msgByte := <-h.broadcast:

			// filtering: check the SessionId and Username
			msg := &wsmsg.WebsocketMessage{}
			if err := proto.Unmarshal(msgByte, msg); err != nil {
				log.Println("fail to unmarshal:", err)
				continue
			}
			fmt.Printf("session %s: user %s\t: %v\n", msg.SessionId, msg.Username, msg.Message)

			for _, client := range h.clients {
				if client.Username == msg.Username {
					if msg.Type != wsmsg.Type_Join {
						log.Printf("not send msg to sender back: username is %v", msg.Username)
						continue
					}
					// before sending welcome message, send all the history to newlyjoined
				}

				if client.SessionID != msg.SessionId {
					log.Printf("send msg only to proper session clients: %s != %s", client.SessionID, msg.SessionId)
					continue
				}

				select {
				case client.send <- msgByte:
				default: // if client cannot consume, clear client
					close(client.send)
					delete(h.clients, client.conn)
				}
			}
		}
	}
}
