package main

import (
	"2do.com/config"
	"2do.com/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"2do.com/handlers"
	"2do.com/repository"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to database
	db := repository.ConnectDB()
	defer db.Close()

	// Initialize repository
	taskRepo := repository.NewPostgresqlClassic(db)

	ctx := context.Background()
	if err := service.Migrate(ctx, taskRepo); err != nil {
		logger.Fatal("Failed to run database migration", zap.Error(err))
	}

	// Initialize service
	//taskService := service.NewTaskService(taskRepo, logger)

	// Initialize router
	router := mux.NewRouter()

	// Register routes
	handlers.RegisterRoutes(router, taskRepo)

	// Start server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in a goroutine
	go func() {
		logger.Info("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
