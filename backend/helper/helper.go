package helper

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
	} else {
		w.Write([]byte(","))
	}
	w.Write(response)
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
