package utils

import (
	"encoding/json"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// Response is the standard envelope returned by every API endpoint.
type Response struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    any               `json:"data"`
	Errors  map[string]string `json:"errors"`
}

// SendSuccess writes a successful JSON response using the standard envelope.
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data any) error {
	return sendJSON(w, statusCode, Response{
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

// SendError writes an error JSON response using the standard envelope.
func SendError(w http.ResponseWriter, statusCode int, message string, errors map[string]string) error {
	return sendJSON(w, statusCode, Response{
		Status:  StatusError,
		Message: message,
		Data:    nil,
		Errors:  errors,
	})
}

func sendJSON(w http.ResponseWriter, statusCode int, response Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}
