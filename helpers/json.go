package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// helper to give response with JSON format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

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

	type errorResponse struct {
		Error string `json:"message"`
	}

	RespondWithJSON(w, code, errorResponse{Error: message})

}