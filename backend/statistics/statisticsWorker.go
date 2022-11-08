package statistics

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Channels struct {
	MsgPerSecondChanel chan CountMsg
}

type CountMsg struct {
	AgentID      uuid.UUID
	MsgPerSecond int
}

type MessageMessagePerSecond struct {
	Type             string `json:"type"`
	MessagePerSecond int    `json:"message_per_second"`
}

func (h *Channels) MessageStats(conn *websocket.Conn) {
	count := 0
	unixtime := time.Now().Unix()
	unixtimeTmp := unixtime
	for el := range h.MsgPerSecondChanel {
		count += el.MsgPerSecond
		if time.Now().Unix() > unixtimeTmp {
			fmt.Println("Messages per second: ", count)
			conn.WriteJSON(MessageMessagePerSecond{Type: "msps", MessagePerSecond: count})
			unixtimeTmp = time.Now().Unix()
			count = 0
		}
	}
}
