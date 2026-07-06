package slots

import (
	"apexcarting/backend/internal/http/response"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filters := ListFilters{
		DateFrom:      q.Get("date_from"),
		DateTo:        q.Get("date_to"),
		TrackConfig:   q.Get("track_config"),
		MarshalID:     q.Get("marshal_id"),
		OnlyAvailable: q.Get("only_available") == "true",
	}

	items, err := h.service.List(r.Context(), filters)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to load slots")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"items": items,
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	slot, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSlotNotFound) {
			response.Error(w, http.StatusNotFound, "slot_not_found", "slot not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "internal_server_error", "failed to load slot")
		return
	}

	response.JSON(w, http.StatusOK, slot)
}
