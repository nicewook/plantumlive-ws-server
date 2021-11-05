// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var debug = flag.Bool("d", false, "display debugging log")

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !(*debug) {
		log.SetOutput(ioutil.Discard)
	}

	hub := newHub()
	go hub.run()

	// router
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("websocket connection received")
		serveWs(hub, w, r)
	})
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("$PORT is not set, use default port %s", port)
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
