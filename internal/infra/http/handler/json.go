package handler

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data any `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return encoder.Encode(data)
}

func jsonResponse(w http.ResponseWriter, status int, data any) error {
	return writeJSON(w, status, &response{data})
}

func jsonErrorResponse(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, &errorResponse{message})
}
