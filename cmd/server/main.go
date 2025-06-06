package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dnakolan/rate-limiter-service/internal/config"
	"github.com/dnakolan/rate-limiter-service/internal/handlers"
	"github.com/dnakolan/rate-limiter-service/internal/services"
	"github.com/dnakolan/rate-limiter-service/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	router := gin.Default()
	gin.SetMode(cfg.Server.GinMode)

	storage := storage.NewRateLimitStorage()
	service := services.NewLimitsService(storage)
	handler := handlers.NewLimitsHandler(service)

	healthHandler := handlers.NewHealthHandler()

	router.GET("/health", healthHandler.GetHealthHandler)

	router.POST("/limits", handler.CreateRateLimitHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	slog.Info("Received terminate, graceful shutdown", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Server exited properly")

	os.Exit(0)
}
