// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	_ "plantumlive-ws-server/proto/wsmsg"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	hub := newHub()
	go hub.run()

	// router
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)

	r.HandleFunc("/ws/{roomID}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		roomID := params["roomID"]
		log.Printf("Room ID is %v", roomID)
		serveWs(hub, roomID, w, r)
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("create new room")
		var roomID string
		serveWs(hub, roomID, w, r)
	})
	http.Handle("/", r)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
