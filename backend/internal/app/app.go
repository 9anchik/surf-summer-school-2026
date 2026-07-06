package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"apexcarting/backend/internal/auth"
	"apexcarting/backend/internal/bookings"
	"apexcarting/backend/internal/config"
	"apexcarting/backend/internal/db"
	httptransport "apexcarting/backend/internal/http"
	"apexcarting/backend/internal/profile"
	"apexcarting/backend/internal/slots"
)

func Run() error {
	cfg := config.Load()

	pool, err := db.NewPostgresPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	slotsRepo := slots.NewRepository(pool)
	slotsService := slots.NewService(slotsRepo)
	slotsHandler := slots.NewHandler(slotsService)

	profileRepo := profile.NewRepository(pool)
	profileService := profile.NewService(profileRepo)
	profileHandler := profile.NewHandler(profileService)

	bookingsRepo := bookings.NewRepository(pool)
	bookingsService := bookings.NewService(bookingsRepo)
	bookingsHandler := bookings.NewHandler(bookingsService)

	router := httptransport.NewRouter(slotsHandler, authHandler, profileHandler, bookingsHandler, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("server started on %s", addr)

	return http.ListenAndServe(addr, router)
}
