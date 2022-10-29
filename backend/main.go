package main

import (
	"encoding/json"
	"fmt"
	"main/agent"
	"main/helper"
	"main/manager"
	"net/http"
	"time"
)

type HandlerChanel struct {
	MsgPerSecondChanel chan int
	IsoStore           manager.IsolateStore
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

func main() {
	runSimulatorInst := &HandlerChanel{
		MsgPerSecondChanel: make(chan int, 300),
		IsoStore:           manager.NewIsolateStore(),
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

func (h *HandlerChanel) runSimulator(w http.ResponseWriter, r *http.Request) {
	defer helper.TimeTrack(time.Now(), "Run Simulator")

	decoder := json.NewDecoder(r.Body)
	var data RunSimulator
	err := decoder.Decode(&data)
	if err != nil {
		helper.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	for i := 0; i < data.Config.Nodes; i++ {
		go agent.New(data.Code)
	}
}
