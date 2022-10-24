package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	v8 "rogchap.com/v8go"
)

func main() {
	http.HandleFunc("/run", runJsCode)
	http.HandleFunc("/ws", socketHandler)

	fs := http.FileServer(http.Dir("./dist"))
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

func runJsCode(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(time.Now(), "runJsCode")
	w.Write([]byte("["))
	defer w.Write([]byte("]"))

	// get json data from post request
	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	var jsCode = data["code"].(string)

	iso := v8.NewIsolate()
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		data, err := interfaceToString(map[string]interface{}{
			"msg":  info.Args()[0].String(),
			"time": time.Now().Format(time.RFC3339),
		})
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return nil
		}

		val, err := v8.NewValue(iso, data)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return nil
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"msg":  info.Args()[0].String(),
			"time": time.Now().Format(time.RFC3339),
		})

		return val
	})
	global := v8.NewObjectTemplate(iso)
	global.Set("test", printfn)
	ctx := v8.NewContext(iso, global)

	_, err = ctx.RunScript(jsCode, "main.js") // execute some JS code
	if err != nil {
		//return err
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
}

func interfaceToString(data interface{}) (string, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
	} else {
		w.Write([]byte(","))
	}
	w.Write(response)
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
