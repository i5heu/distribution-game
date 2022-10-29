package agent

import (
	"fmt"
	"reflect"
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
	Chanel chan Data
}

type agentCtx struct {
	info   Info
	agents *map[uuid.UUID]Agent
}

func New(code string, agentData Agent, agents *map[uuid.UUID]Agent, seeds []uuid.UUID) {
	info := Info{
		AgentID: agentData.NodeID,
		Seeds:   seeds,
	}
	ctx := agentCtx{
		info:   info,
		agents: agents,
	}

	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(
		map[string]map[string]reflect.Value{
			"custom/custom": {
				"Data":    reflect.ValueOf((*Data)(nil)),
				"Info":    reflect.ValueOf((*Info)(nil)),
				"GetInfo": reflect.ValueOf(info.nodeInfo),
				"Send":    reflect.ValueOf(ctx.nodeSend),
			},
		})

	_, err := i.Eval(code)
	if err != nil {
		fmt.Println("CompilerError:", err)
	}

	agentFuncReceive, err := i.Eval("Receive")
	if err != nil {
		fmt.Println(err)
	}
	receive := agentFuncReceive.Interface().(func(Data))
	go ReceiveWorker(agentData.Chanel, receive)

	agentFuncInit, err := i.Eval("Init")
	if err != nil {
		fmt.Println(err)
	}
	init := agentFuncInit.Interface().(func())
	init()

}

func ReceiveWorker(dataChan chan Data, receive func(Data)) {
	count := 0
	unixTime := time.Now().Unix()

	for data := range dataChan {
		count++

		if unixTime < time.Now().Unix() {
			fmt.Println("ReceiveWorker", count)
			count = 0
			unixTime = time.Now().Unix()
		}

		receive(data)
	}
}

func (i Info) nodeInfo() Info {
	return Info{
		AgentID: i.AgentID,
		Seeds:   i.Seeds,
	}
}

func (ctx agentCtx) nodeSend(target uuid.UUID, msg string) {
	(*ctx.agents)[target].Chanel <- Data{
		SenderID: ctx.info.AgentID,
		Msg:      msg,
	}
}
