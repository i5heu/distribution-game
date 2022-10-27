package main

import (
	"encoding/json"
	"fmt"
	"main/agent"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	v8 "rogchap.com/v8go"
)

type HandlerChanel struct {
	MsgPerSecondChanel chan int
	//whatever
}

func main() {

	runSimulatorInst := &HandlerChanel{
		MsgPerSecondChanel: make(chan int, 300),
	}

	http.HandleFunc("/run", runSimulatorInst.runSimulator)
	http.HandleFunc("/ws", runSimulatorInst.socketHandler)

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

func (h *HandlerChanel) runSimulator(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data RunSimulator
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			fmt.Println("Recovered", err)
		}
	}()

	defer timeTrack(time.Now(), "Run Simulator")

	var jsCode = data.Code
	processors := runtime.NumCPU()

	// create a new Isolate with the number of available CPUs
	isos := make([]*v8.Isolate, processors)
	for i := 0; i < processors; i++ {
		isos[i] = v8.NewIsolate()
	}

	indexIsos := 0
	for i := 0; i < data.Config.Nodes; i++ {
		agent.New(uuid.New().String(), jsCode, isos[indexIsos], h.MsgPerSecondChanel)
		indexIsos++
		if indexIsos == processors {
			indexIsos = 0
		}
	}
}

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
