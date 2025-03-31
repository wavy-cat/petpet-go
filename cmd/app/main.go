package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Get config
	cfg, err := config.GetConfig("config.yml")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Create a cache object
	var cachePNG, cacheGIF cache.BytesCache

	switch cfg.Storage {
	case "memory":
		cacheGIF, err = memory.NewLRUCache(cfg.MemoryCapacity)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}

		cachePNG, err = memory.NewLRUCache(cfg.MemoryCapacity)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}
	case "fs":
		cacheGIF, err = fs.NewFileSystemCache(cfg.FSPath)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}

		cachePNG, err = fs.NewFileSystemCache(cfg.FSPath)
		if err != nil {
			logger.Fatal("Error creating memory cache object", zap.Error(err))
		}
	case "":
	default:
		logger.Warn("Passed an incorrect storage type for the cache. The cache will be disabled")
	}

	// Add proxy
	var transport *http.Transport

	if cfg.URL != "" {
		proxyURL, err := url.Parse(cfg.URL)
		if err != nil {
			logger.Warn("Failed to parse proxy URL. The server will be started without using a proxy",
				zap.Error(err))
		} else {
			transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}

	// Create a bot object
	discordBot := discord.NewBot(cfg.BotToken)

	// Setting up the service
	providers := map[string]repository.AvatarProvider{
		"discord": repository.NewDiscordAvatarProvider(discordBot),
	}
	gifService := service.NewGIFService(cacheGIF, providers, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
	apngService := service.NewAPngService(cachePNG, providers, petpet.DefaultConfig)

	// Set up routing
	router := mux.NewRouter()

	gifHandler := &middleware.Logging{
		Logger: logger,
		Next:   ds_gif.NewHandler(gifService, transport),
	}
	router.Handle("/ds/{user_id}.gif", gifHandler).Methods(http.MethodGet)
	router.Handle("/ds/{user_id}", gifHandler).Methods(http.MethodGet)

	apngHandler := &middleware.Logging{
		Logger: logger,
		Next:   ds_apng.NewHandler(apngService, transport),
	}
	router.Handle("/ds/{user_id}.apng", apngHandler).Methods(http.MethodGet)

	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("See documentation on GitHub: https://github.com/wavy-cat/petpet-go"))
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
	}).Methods(http.MethodGet, http.MethodHead)

	// Set up the server
	var serverAddr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	go func() {
		logger.Info("Starting the HTTP server...", zap.String("Address", serverAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed:", zap.Error(err))
		}
	}()

	// Waiting for completion signal
	<-stop

	// Create a context with a timeout to shut down the server gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Millisecond)
	defer cancel()

	logger.Info("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
