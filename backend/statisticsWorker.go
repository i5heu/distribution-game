package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func (h *HandlerChanel) messageCounter(conn *websocket.Conn) {
	count := 0
	unixtime := time.Now().Unix()
	unixtimeTmp := unixtime
	for el := range h.MsgPerSecondChanel {
		count += el
		if time.Now().Unix() > unixtimeTmp {
			fmt.Println("Messages per second: ", count)
			conn.WriteJSON(MessageMessagePerSecond{Type: "msps", MessagePerSecond: count})
			unixtimeTmp = time.Now().Unix()
			count = 0
		}
	}
}
