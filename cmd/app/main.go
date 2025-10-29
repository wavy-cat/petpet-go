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

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/ds_gif"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/repository/avatar/discord"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/cache/fs"
	"github.com/wavy-cat/petpet-go/pkg/cache/memory"
	"github.com/wavy-cat/petpet-go/pkg/cache/s3"
	"github.com/wavy-cat/petpet-go/pkg/logger-presets/gcp"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"go.uber.org/zap"
)

const serviceName = "petpet"

func main() {
	// Get config
	cfg, err := config.GetConfig("config.yml")
	if err != nil {
		panic(err)
	}

	// Setting up a logger-presets
	var logger *zap.Logger

	switch cfg.Logger.Preset {
	case config.ProdPreset:
		logger, err = zap.NewProduction()
	case config.DevPreset:
		logger, err = zap.NewDevelopment()
	case config.GCPPreset:
		logger, err = gcp.NewGCPLogger()
	default:
		fmt.Println("Logging is disabled by default. To enable it, select a logger preset in the configuration.")
		logger = zap.NewNop()
	}

	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Error syncing logger-presets:", err)
		}
	}(logger)

	// Create a cache instance
	var cacheInstance cache.BytesCache

	switch cfg.Storage {
	case "memory":
		cacheInstance, err = memory.NewLRUCache(cfg.Memory.Capacity)
		if err != nil {
			logger.Fatal("Error creating memory cacheInstance object", zap.Error(err))
		}
	case "fs":
		cacheInstance, err = fs.NewFileSystemCache(cfg.FS.Path)
		if err != nil {
			logger.Fatal("Error creating filesystem cacheInstance object", zap.Error(err))
		}
	case "s3":
		if cfg.S3.Bucket == "" {
			logger.Fatal("S3 bucket name is required for S3 cacheInstance")
		}

		cacheInstance, err = s3.NewS3Cache(cfg.S3.Bucket, cfg.S3.Endpoint, cfg.S3.Region, cfg.S3.AccessKey, cfg.S3.SecretKey)
		if err != nil {
			logger.Fatal("Error creating S3 cacheInstance object", zap.Error(err))
		}
	case "":
	default:
		logger.Warn("Passed an incorrect storage type for the cacheInstance. The cacheInstance will be disabled")
	}

	// Add proxy
	var transport *http.Transport

	if cfg.URL != "" {
		proxyURL, err := url.Parse(cfg.URL)
		if err != nil {
			logger.Fatal("Failed to parse proxy URL.", zap.Error(err))
		} else {
			transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}

	// Set up routing
	r := chi.NewRouter()

	r.Use(middleware.Logger(logger, serviceName))

	if cfg.Heartbeat.Enable {
		r.Use(chiMiddleware.Heartbeat(cfg.Heartbeat.Path))
	}

	if cfg.Throttle.Enable {
		r.Use(chiMiddleware.ThrottleBacklog(
			cfg.Throttle.Limit,
			cfg.Throttle.Backlog,
			time.Duration(cfg.Throttle.BacklogTimeout)*time.Second))
	}

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("See documentation on GitHub: https://github.com/wavy-cat/petpet-go"))
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
	})

	// Add Discord service
	gifService := service.NewGIFService(cacheInstance,
		discord.NewDiscordAvatarProvider(cfg.BotToken),
		petpet.DefaultConfig,
		quantizers.HierarhicalQuantizer{})

	gifHandler := ds_gif.NewHandler(gifService, transport)
	r.Method(http.MethodGet, "/ds/{user_id}.gif", gifHandler)
	r.Method(http.MethodGet, "/ds/{user_id}", gifHandler)

	// Set up the server
	var serverAddr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
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

	logger.Info("Closing the cache...")
	if err := cacheInstance.Close(); err != nil {
		logger.Error("Error closing cache", zap.Error(err), zap.String("cache_type", cfg.Storage))
	}

	logger.Info("Server exited properly")
}
