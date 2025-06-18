package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := healthResponse{
		Status: http.StatusText(http.StatusOK),
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
