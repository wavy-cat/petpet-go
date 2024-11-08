package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/ds_apng"
	"github.com/wavy-cat/petpet-go/internal/handler/http/ds_gif"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/repository"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/cache/fs"
	"github.com/wavy-cat/petpet-go/pkg/cache/memory"
	"github.com/wavy-cat/petpet-go/pkg/discord"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Setting up a logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Error syncing logger:", err)
		}
	}(logger)

	// Create a cache object
	var cachePNG, cacheGIF cache.BytesCache

	switch config.CacheStorage {
	case "memory":
		cacheGIF, err = memory.NewLRUCache(config.CacheMemoryCapacity)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}

		cachePNG, err = memory.NewLRUCache(config.CacheMemoryCapacity)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}
	case "fs":
		cacheGIF, err = fs.NewFileSystemCache(config.CacheFSPath)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}

		cachePNG, err = fs.NewFileSystemCache(config.CacheFSPath)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}
	case "":
	default:
		logger.Warn("Passed an incorrect storage type for the cache. The cache will be disabled")
	}

	// Create a bot object
	discordBot := discord.NewBot(config.BotToken)

	// Setting up the service
	providers := map[string]repository.AvatarProvider{
		"discord": repository.NewDiscordAvatarProvider(discordBot),
	}
	gifService := service.NewGIFService(cacheGIF, providers, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
	apngService := service.NewAPngService(cachePNG, providers, petpet.DefaultConfig)

	// Set up routing
	router := mux.NewRouter()

	gifHandle := middleware.Logging{
		Logger: logger,
		Next:   ds_gif.NewHandler(gifService),
	}
	router.Handle("/ds/{user_id}.gif", &gifHandle).Methods(http.MethodGet)

	apngHandle := middleware.Logging{
		Logger: logger,
		Next:   ds_apng.NewHandler(apngService),
	}
	router.Handle("/ds/{user_id}.apng", &apngHandle).Methods(http.MethodGet)

	router.Handle("/ds/{user_id}", &gifHandle).Methods(http.MethodGet)

	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("Waiting for something to happen?"))
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
	}).Methods(http.MethodGet, http.MethodHead)

	// Set up the server
	srv := &http.Server{
		Addr:    config.HTTPAddress,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	go func() {
		logger.Info("Starting the HTTP server...", zap.String("Address", config.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed:", zap.Error(err))
		}
	}()

	// Waiting for completion signal
	<-stop

	// Create a context with a timeout to shut down the server gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ShutdownTimeout)*time.Millisecond)
	defer cancel()

	logger.Info("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
