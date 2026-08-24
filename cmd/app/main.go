package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wavy-cat/petpet-go/internal/cachefactory"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/loggerfactory"
	"github.com/wavy-cat/petpet-go/internal/server"
	"go.uber.org/zap"
)

const serviceName = "petpet"

func main() {
	// Get config
	cfg, err := config.GetConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// Setting up a logger-presets
	logger, err := loggerfactory.NewWithService(cfg.Logger, serviceName)
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	// Create a cache instance
	cacheInstance, err := cachefactory.New(cfg.Cache)
	if err != nil {
		logger.Fatal("Error creating cache", zap.Error(err))
	}
	defer func() {
		logger.Info("Closing the cache...")
		if err := cacheInstance.Close(); err != nil {
			logger.Error("Error closing cache", zap.Error(err))
		}
	}()

	// Set up routing
	server, err := server.NewServer(cfg, logger, cacheInstance)
	if err != nil {
		logger.Fatal("Error when creating server", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	go func() {
		logger.Info("Starting the HTTP server...")

		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Waiting for completion signal
	<-stop

	// Create a context with a timeout to shut down the server gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	logger.Info("Shutting down the server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
