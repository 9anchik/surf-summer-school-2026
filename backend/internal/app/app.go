package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"apexcarting/backend/internal/config"
	"apexcarting/backend/internal/db"
	httptransport "apexcarting/backend/internal/http"
)

func Run() error {
	cfg := config.Load()

	pool, err := db.NewPostgresPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	router := httptransport.NewRouter()

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("server started on %s", addr)

	return http.ListenAndServe(addr, router)
}
