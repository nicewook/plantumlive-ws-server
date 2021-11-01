// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
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

const (
	TypeConnected = "connected"
	TypeJoin      = "join-session"
	TypeMsg       = "message"
)

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			welcomeMsg := &wsmsg.WebsocketMessage{
				Type:    TypeConnected,
				Message: "You are connected to websocket server",
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
			receivedMsg := &wsmsg.WebsocketMessage{}
			log.Printf("%+v", string(msgByte))
			if err := proto.Unmarshal(msgByte, receivedMsg); err != nil {
				log.Println("fail to unmarshal:", err)
				continue
			}

			for client := range h.clients {
				if client.SessionID != receivedMsg.SessionId {
					log.Printf("not send msg: Client RoomID: %v, receivedMsg RoomID: %v", client.SessionID, receivedMsg.SessionId)
					continue
				}
				if client.Username == receivedMsg.Username {
					log.Printf("not send msg to sender back: username is %v", receivedMsg.Username)
					continue
				}
				select {
				case client.send <- msgByte:
				default: // if client cannot consume, clear client
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
