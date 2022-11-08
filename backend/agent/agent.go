package agent

import (
	"context"
	"fmt"
	"main/helper"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Data struct {
	SenderID uuid.UUID
	Msg      string `json:"msg"`
}

type Info struct {
	AgentID uuid.UUID
	Seeds   []uuid.UUID
}

type Agent struct {
	NodeID uuid.UUID `json:"node_id"`
	Chanel *chan Data
}

type agentEnv struct {
	info   Info
	agents *map[uuid.UUID]Agent
}

func New(code string, agentData Agent, agents *map[uuid.UUID]Agent, seeds []uuid.UUID, ctx context.Context) {
	info := Info{
		AgentID: agentData.NodeID,
		Seeds:   seeds,
	}
	agentEnv := agentEnv{
		info:   info,
		agents: agents,
	}

	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(
		map[string]map[string]reflect.Value{
			"dg/dg": {
				"Data":    reflect.ValueOf((*Data)(nil)),
				"Info":    reflect.ValueOf((*Info)(nil)),
				"GetInfo": reflect.ValueOf(info.nodeInfo),
				"Send":    reflect.ValueOf(agentEnv.nodeSend),
			},
		})

	_, err := i.EvalWithContext(ctx, code)
	if err != nil {
		fmt.Println("CompilerError:", err)
	}

	agentFuncReceive, err := i.EvalWithContext(ctx, "Receive")
	if err != nil {
		fmt.Println(err)
	}
	receive := agentFuncReceive.Interface().(func(Data))
	go ReceiveWorker(*agentData.Chanel, receive, *agents)

	agentFuncClose, err := i.Eval("Close")
	if err != nil {
		fmt.Println(err)
	}
	close := agentFuncClose.Interface().(func())

	agentFuncInit, err := i.EvalWithContext(ctx, "Init")
	if err != nil {
		fmt.Println(err)
	}
	init := agentFuncInit.Interface().(func())
	go init()

	<-ctx.Done()
	close()
}

func ReceiveWorker(dataChan chan Data, receive func(Data), agents map[uuid.UUID]Agent) {
	count := 0
	// unixTime := time.Now().Unix()

	missingAgents := make(map[uuid.UUID]int)
	// fill missingAgents
	for _, agent := range agents {
		missingAgents[agent.NodeID] = 0
	}

	for data := range dataChan {
		count++

		msg, _ := strconv.ParseInt(data.Msg, 10, 64)

		t := time.Unix(0, msg)

		helper.TimeTrack(t, "Run Receive "+data.SenderID.String())

		// if unixTime < time.Now().Unix() {
		// 	time.Sleep(time.Second / 8)

		// 	fmt.Println("\n\n\n", "ReceiveWorker", count, "messages in queue", len(dataChan))
		// 	count = 0
		// 	unixTime = time.Now().Unix()

		// 	//check if all agents are still alive
		// 	for _, agent := range agents {
		// 		fmt.Println("Agent", agent.NodeID, " posted ", missingAgents[agent.NodeID], " messages")
		// 	}

		// 	fmt.Println("\n\n\n")

		// 	//reset missingAgents
		// 	for _, agent := range agents {
		// 		missingAgents[agent.NodeID] = 0
		// 	}
		// }

		receive(data)
	}
}

func (i Info) nodeInfo() Info {
	return Info{
		AgentID: i.AgentID,
		Seeds:   i.Seeds,
	}
}

func (ae agentEnv) nodeSend(target uuid.UUID, msg string) {
	defer helper.TimeTrack(time.Now(), "Run Target "+target.String())

	time.Sleep(time.Second / 10)

	targetAgent := (*ae.agents)[target]

	if IsOpen(*targetAgent.Chanel) {
		*targetAgent.Chanel <- Data{
			SenderID: ae.info.AgentID,
			Msg:      msg,
		}
	}
}

func IsOpen(ch <-chan Data) bool {
	select {
	case <-ch:
		return false
	default:
	}

	return true
}
