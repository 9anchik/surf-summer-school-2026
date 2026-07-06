package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string         `json:"code,omitempty"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func ErrorWithDetails(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
	details map[string]any,
) {
	JSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}
