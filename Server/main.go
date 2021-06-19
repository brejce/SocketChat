// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8989", "http service address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/message", func(rw http.ResponseWriter, r *http.Request) {
		message(hub, rw, r)
	})
	http.HandleFunc("/register", func(rw http.ResponseWriter, r *http.Request) {
		Register(rw, r)
	})
	http.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		Login(rw, r)
	})
	http.HandleFunc("/dislogin", func(rw http.ResponseWriter, r *http.Request) {
		DisLogin(rw, r)
	})
	http.HandleFunc("/chanagepasswd", func(rw http.ResponseWriter, r *http.Request) {
		ChanagePasswd(rw, r)
	})
	http.HandleFunc("/getallmessage", func(rw http.ResponseWriter, r *http.Request) {
		GetAllMsg(rw, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
