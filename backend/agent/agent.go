package agent

import (
	"fmt"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Data struct {
	Msg string `json:"msg"`
}

func New(code string) {

	fmt.Println("111")
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(
		map[string]map[string]reflect.Value{
			"custom/custom": {
				"Data": reflect.ValueOf((*Data)(nil)),
				"Send": reflect.ValueOf(nodeSend),
			},
		})

	_, err := i.Eval(code)
	if err != nil {
		fmt.Println("ERRRROLL", err)
	}

	fmt.Println("Done xxxxxxxx Done")

	v, err := i.Eval("Receive")
	if err != nil {
		fmt.Println(err)
	}

	bar := v.Interface().(func(Data))

	bar(Data{Msg: "Hello World"})
}

func nodeSend(data interface{}) {
	fmt.Println("nodeSend", data)
}

func interfaceToString(data Data) {
	fmt.Println("interfaceToString", data)
}
