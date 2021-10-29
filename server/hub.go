// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"log"
	wsmsg "plantumlive-ws-server/wsmsg"

	"google.golang.org/protobuf/proto"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type WebsocketMessage struct {
	RoomID   string `json:"roomid,omitempty"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message,omitempty"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// make msg and serialize and send to client welcom
			// welcomeMsg := WebsocketMessage{
			// 	RoomID:   client.RoomID,
			// 	Username: client.Username,
			// 	Message:  "welcome",
			// }
			welcomeMsg := &wsmsg.WebsocketMessage{
				Type:      "welcome",
				SessionId: client.RoomID,
				Username:  client.Username,
				Message:   "welcome",
			}

			msgBytes, err := proto.Marshal(welcomeMsg)
			if err != nil {
				log.Println("fail to marshal:", err)
				continue
			}
			client.send <- msgBytes

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msgByte := <-h.broadcast:
			// byte to struct
			// check the room id and sender name - filtering
			var receivedMsg WebsocketMessage
			if err := json.Unmarshal(msgByte, &receivedMsg); err != nil {
				log.Println("fail to marshal:", err)
				continue
			}

			for client := range h.clients {
				// temporary disable
				// if client.RoomID != receivedMsg.RoomID {
				// 	log.Printf("not send msg: Client RoomID: %v, receivedMsg RoomID: %v", client.RoomID, receivedMsg.RoomID)
				// 	continue
				// }
				if client.Username == receivedMsg.Username {
					log.Printf("not send msg to sender back: username is %v", receivedMsg.Username)
					continue
				}
				select {
				case client.send <- msgByte:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
