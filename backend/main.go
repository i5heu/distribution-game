package main

import (
	"encoding/json"
	"fmt"
	"main/agent"
	"main/helper"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func main() {
	http.HandleFunc("/run", runSimulator)
	http.HandleFunc("/ws", socketHandler)

	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	} else {
		fmt.Println("Server started on port 3333")
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

type Config struct {
	Nodes      int `json:"nodes"`
	MsgsSNode  int `json:"msgs_s_node"`
	Datasets   int `json:"datasets"`
	DatasetsS  int `json:"datasets_s"`
	Seeds      int `json:"seeds"`
	Iterations int `json:"iterations"`
	Timeout    int `json:"timeout"`
}

type RunSimulator struct {
	Config Config `json:"config"`
	Code   string `json:"code"`
}

func runSimulator(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "Run Simulator")
	w.Write([]byte("["))
	defer w.Write([]byte("]"))

	// get json data from post request
	decoder := json.NewDecoder(r.Body)
	var data RunSimulator
	err := decoder.Decode(&data)
	if err != nil {
		helper.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	var jsCode = data.Code

	for i := 0; i < data.Config.Nodes; i++ {
		go agent.New(uuid.New().String(), jsCode)
	}
}

type Message struct {
	ThreadID    uuid.UUID       `json:"thread_id"`
	MessageType string          `json:"type"`
	Data        json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

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
