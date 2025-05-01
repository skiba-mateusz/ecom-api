package http

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type response struct {
	Data any `json:"data"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
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
