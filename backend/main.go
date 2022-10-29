package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/agent"
	"main/helper"
	"main/manager"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
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
	var err error
	// get the file descriptor for the log file
	// f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	redirectStderr(os.Stdout)

	runSimulatorInst := &HandlerChanel{
		MsgPerSecondChanel: make(chan int, 300),
		IsoStore:           manager.NewIsolateStore(),
	}

	http.HandleFunc("/run", runSimulatorInst.runSimulator)
	http.HandleFunc("/ws", runSimulatorInst.socketHandler)

	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)

	err = http.ListenAndServe(":3333", nil)
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

	agents := make(map[uuid.UUID]agent.Agent)

	var simpleSeed uuid.UUID

	for i := 0; i < data.Config.Nodes; i++ {
		agentId := uuid.New()
		agents[agentId] = agent.Agent{
			NodeID: agentId,
			Chanel: make(chan agent.Data, 100),
		}

		if simpleSeed == uuid.Nil {
			simpleSeed = agentId
		}
	}

	for _, agentData := range agents {
		go agent.New(data.Code, agentData, &agents, []uuid.UUID{agents[simpleSeed].NodeID})
	}

	// close channels after timeout
	go func() {
		time.Sleep(time.Duration(data.Config.Timeout) * time.Second)
		fmt.Println("Timeout... closing channels")
		for _, agentData := range agents {
			close(agentData.Chanel)
		}
	}()
}

func redirectStderr(f *os.File) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
