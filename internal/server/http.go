package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/trace"

	"contract_ease/internal/config"
	"contract_ease/internal/database"
)

func BuildHTTPServer(ctx context.Context, tp trace.TracerProvider) *http.Server {
	cfg := config.LoadConfig()
	port, _ := strconv.Atoi(cfg.App.Port)
	db := database.New(ctx, cfg)

	app := NewServer(db, port)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      app.RegisterRoutes(tp),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
