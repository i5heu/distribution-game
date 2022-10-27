package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type MessageMessagePerSecond struct {
	Type             string `json:"type"`
	MessagePerSecond int    `json:"message_per_second"`
}
type Message struct {
	ThreadID    uuid.UUID       `json:"thread_id"`
	MessageType string          `json:"type"`
	Data        json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{}

func (h *HandlerChanel) socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	go h.messageCounter(conn)
	go h.IsoStore.IsolateManager()

	messageCount := -1

	// The event loop
	for {
		messageCount++
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error().AnErr("Error during message reading", err)
			break
		}

		data := Message{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Error().Err(err).Msg("Error during message unmarshalling")
			conn.WriteMessage(messageType, []byte("Error during message unmarshalling"))
			break
		}

		conn.WriteMessage(messageType, []byte("{\"data\":\"Working\"}"))
	}
}
