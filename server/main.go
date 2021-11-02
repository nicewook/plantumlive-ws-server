// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var port = flag.String("port", ":8080", "http service address. default is :8080")
var debug = flag.Bool("debug", false, "display debugging log")

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

	fmt.Printf("websocket server started on http://localhost%s\n", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
