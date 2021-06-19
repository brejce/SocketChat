// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, message, err := c.conn.ReadMessage() //message代表无限可能,loginregister等功能都可以在这里做
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg Message
		e := json.Unmarshal(message, &msg)
		CheckError(e)
		if !IsOurUser(msg.Name) {
			SendToClient(c.conn, mt, "You are not Login Users")
			break
		}
		SetMessage(msg)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 将连接升级为 webSocket 并将其加入消息队列
func message(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()

}
func Register(rw http.ResponseWriter, r *http.Request) {
	s := ""
	data := User{}
	json.NewDecoder(r.Body).Decode(&data)
	b := SetUser(data)
	if b { //成功保存此用户到数据库,用户注册成功
		s = "yes"
	} else { //已存在,用户注册失败
		s = "no"
	}
	rw.Write([]byte(s))
}
func DisLogin(rw http.ResponseWriter, r *http.Request) {
	data := User{}
	json.NewDecoder(r.Body).Decode(&data)
	DeletLoginUser(data)
	rw.Write([]byte("yes"))
}
func Login(rw http.ResponseWriter, r *http.Request) {
	s := "defalut"
	data := User{}
	json.NewDecoder(r.Body).Decode(&data)
	user, b := GetUser(data)
	if b { //用户名正确
		if data.Passwd == user.Passwd { //密码正确,加入登录列表,登录成功
			s = "yes"
			if !IsOurUser(data.Name) {
				u1 := user
				u1.Name = data.Name //只保存用户名
				var l = len(UserMap)
				UserMap[l+1] = u1
			}
		} else { //密码错误,登录失败
			s = "no"
		}
	} else { //用户名错误
		s = "no"
	}
	rw.Write([]byte(s))
}
func ChanagePasswd(rw http.ResponseWriter, r *http.Request) {
	s := "defalut"
	data := User{}
	json.NewDecoder(r.Body).Decode(&data)
	if SetPasswd(data) {
		//修改成功
		s = "yes"
	} else {
		//修改失败,检查账号和密码,如还是失败联系管理员
		s = "no"
	}
	rw.Write([]byte(s))
}
func GetAllMsg(rw http.ResponseWriter, r *http.Request) {
	//  IsOurUser()
	s := GetAllMessageRange()
	data, e := json.Marshal(s)
	CheckError(e)
	rw.Write(data)
}

//将websocket.Conn.WriteMessage封装一下
func SendToClient(ws *websocket.Conn, mt int, s string) {
	err := ws.WriteMessage(mt, []byte(s))
	if err != nil {
		return
	}
}
