package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
}

// helper to give response with JSON format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(Response{Status: "SUCCESS", Data: payload})

	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

func RespondWithErr(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("responding with 5xx err: %v", message)
	}

	response, err := json.Marshal(Response{Status: "ERROR", Error: message})

	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}
