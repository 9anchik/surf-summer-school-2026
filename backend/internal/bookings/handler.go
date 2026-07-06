package bookings

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	authmiddleware "apexcarting/backend/internal/http/middleware"
	"apexcarting/backend/internal/http/response"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	booking, err := h.service.Create(r.Context(), userID, req, idempotencyKey)
	if err != nil {
		switch {
		case errors.Is(err, ErrMissingIdempotencyKey):
			response.Error(w, http.StatusBadRequest, "missing_idempotency_key", "Idempotency-Key header is required")
		case errors.Is(err, ErrInvalidRequest):
			response.Error(w, http.StatusBadRequest, "invalid_booking_request", "invalid booking request")
		case errors.Is(err, ErrSlotNotFound):
			response.Error(w, http.StatusNotFound, "slot_not_found", "slot not found")
		case errors.Is(err, ErrSlotUnavailable):
			response.Error(w, http.StatusGone, "slot_cancelled", "slot is unavailable")
		case errors.Is(err, ErrSlotFull):
			response.Error(w, http.StatusConflict, "slot_full", "not enough free seats")
		case errors.Is(err, ErrRentalUnavailable):
			response.Error(w, http.StatusConflict, "rental_unavailable", "not enough rental equipment")
		default:
			response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to create booking")
		}
		return
	}

	response.JSON(w, http.StatusCreated, booking)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)

	items, err := h.service.ListByUser(r.Context(), userID, ListBookingsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to load bookings")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"limit":  limit,
		"offset": offset,
	})
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	bookingID := chi.URLParam(r, "id")

	booking, err := h.service.Cancel(r.Context(), userID, bookingID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			response.Error(w, http.StatusBadRequest, "invalid_booking_request", "invalid booking id")
		case errors.Is(err, ErrBookingNotFound):
			response.Error(w, http.StatusNotFound, "booking_not_found", "booking not found")
		case errors.Is(err, ErrBookingNotActive):
			response.Error(w, http.StatusConflict, "already_cancelled", "booking is not active")
		case errors.Is(err, ErrSlotAlreadyStarted):
			response.Error(w, http.StatusUnprocessableEntity, "slot_started", "slot already started")
		default:
			response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to cancel booking")
		}
		return
	}

	response.JSON(w, http.StatusOK, booking)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	bookingID := chi.URLParam(r, "id")

	booking, err := h.service.GetByID(r.Context(), userID, bookingID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRequest):
			response.Error(w, http.StatusBadRequest, "invalid_booking_request", "invalid booking id")
		case errors.Is(err, ErrBookingNotFound):
			response.Error(w, http.StatusNotFound, "booking_not_found", "booking not found")
		default:
			response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to load booking")
		}
		return
	}

	response.JSON(w, http.StatusOK, booking)
}
