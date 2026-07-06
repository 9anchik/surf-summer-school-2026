package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	authmiddleware "apexcarting/backend/internal/http/middleware"
	"apexcarting/backend/internal/http/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	user, err := h.service.Get(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			response.Error(w, http.StatusNotFound, "profile_not_found", "profile not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to load profile")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	var req UpdateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	user, err := h.service.Update(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, ErrInvalidProfileData) {
			response.Error(w, http.StatusBadRequest, "invalid_profile_data", "invalid profile data")
			return
		}

		if errors.Is(err, ErrProfileNotFound) {
			response.Error(w, http.StatusNotFound, "profile_not_found", "profile not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to update profile")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	err := h.service.Delete(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			response.Error(w, http.StatusNotFound, "profile_not_found", "profile not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal_error", "failed to delete profile")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "profile deleted",
	})
}
