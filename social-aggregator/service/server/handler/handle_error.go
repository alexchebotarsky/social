package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func handleError(ctx context.Context, w http.ResponseWriter, err error, statusCode int, shouldLog bool) {
	if shouldLog {
		log.Printf("Handler error: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(errorResponse{
		Error:      fmt.Sprintf("%s: %d", http.StatusText(statusCode), statusCode),
		StatusCode: statusCode,
	})
	if err != nil {
		log.Printf("Error encoding error response to JSON: %v", err)
	}
}
