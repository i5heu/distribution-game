package agent

import (
	"encoding/json"
	"fmt"

	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

func New(module *wasmer.Module, store *wasmer.Store) {

	wasiEnv, _ := wasmer.NewWasiStateBuilder("wasi-program").
		// Choose according to your actual situation
		// Argument("--foo").
		// Environment("ABC", "DEF").
		// MapDirectory("./", ".").
		Finalize()
	importObject, err := wasiEnv.GenerateImportObject(store, module)

	instance, err := wasmer.NewInstance(module, importObject)

	if err != nil {
		panic(fmt.Sprintln("Failed to instantiate the module:", err))
	}

	addOne, err := instance.Exports.GetFunction("add_one")

	if err != nil {
		panic(fmt.Sprintln("Failed to get the `add_one` function:", err))
	}

	// // repeat a 1000 times
	for i := 0; i < 100; i++ {
		values, err := addOne(2)

		if err != nil {
			panic(fmt.Sprintln("Failed to call the `add_one` function:", err))
		}

		if values.(int32) != 4 {
			fmt.Println("Expected 4, got", values)
		}
	}

	instance.Close()
}

func interfaceToString(data interface{}) (string, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
