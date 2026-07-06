package auth

import (
	"apexcarting/backend/internal/http/response"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req SendOTPRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	if req.Phone == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "phone is required")
		return
	}

	code, err := h.service.SendOTP(r.Context(), req.Phone)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to send otp")
		return
	}

	response.JSON(w, http.StatusOK, SendOTPResponse{
		Message: "otp sent",
		Code:    code,
	})
}

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req VerifyOTPRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	if req.Phone == "" || req.Code == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "phone and code are required")
		return
	}

	resp, err := h.service.VerifyOTP(r.Context(), req.Phone, req.Code, req.Name)
	if err != nil {
		if errors.Is(err, ErrInvalidOTP) {
			response.Error(w, http.StatusUnauthorized, "invalid_otp", "invalid otp code")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to verify otp")
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{
		"message": "logged out",
	})
}
