package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse writes a standard JSON success envelope with the provided
// status and data payload.
func SuccessResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// ErrorResponse writes a standard JSON error envelope with the provided status
// and error message.
func ErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
