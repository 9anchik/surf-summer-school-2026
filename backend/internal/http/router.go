package http

import (
	"net/http"

	"apexcarting/backend/internal/auth"
	"apexcarting/backend/internal/bookings"
	authmiddleware "apexcarting/backend/internal/http/middleware"
	"apexcarting/backend/internal/profile"
	"apexcarting/backend/internal/slots"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	slotsHandler *slots.Handler,
	authHandler *auth.Handler,
	profileHandler *profile.Handler,
	bookingsHandler *bookings.Handler,
	jwtSecret string,
) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/otp/send", authHandler.SendOTP)
		r.Post("/auth/otp/verify", authHandler.VerifyOTP)

		r.Get("/slots", slotsHandler.List)
		r.Get("/slots/{id}", slotsHandler.GetByID)

		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.Auth(jwtSecret))

			r.Post("/auth/logout", authHandler.Logout)

			r.Get("/profile", profileHandler.Get)
			r.Patch("/profile", profileHandler.Update)
			r.Delete("/profile", profileHandler.Delete)

			r.Get("/bookings", bookingsHandler.List)
			r.Post("/bookings", bookingsHandler.Create)
			r.Get("/bookings/{id}", bookingsHandler.GetByID)
			r.Post("/bookings/{id}/cancel", bookingsHandler.Cancel)
		})
	})

	return r
}
