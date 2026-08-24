package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto/tls"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/wavy-cat/petpet-go/internal/cachefactory"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/custom"
	discord_handler "github.com/wavy-cat/petpet-go/internal/handler/http/discord"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/avatarproviders/discord"
	"github.com/wavy-cat/petpet-go/pkg/logger-presets/gcp"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
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
	var logger *zap.Logger

	switch cfg.Preset {
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
		_ = logger.Sync()
	}(logger)

	// Create a cache instance
	if cfg.Storage == "" {
		logger.Info("The storage type is not specified. Caching will be disabled")
	}

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

	if cfg.Discord.Enable {
		if cfg.BotToken == "" {
			logger.Fatal("Discord bot token is required when Discord is enabled")
		}

		discordGifService := service.NewGIFService(cacheInstance,
			discord.NewProvider(cfg.BotToken),
			petpet.DefaultConfig,
			quantizers.HierarhicalQuantizer{})
		gifHandler := discord_handler.NewHandler(discordGifService)

		r.Get("/discord/{user_id}.gif", gifHandler)
		r.Get("/discord/{user_id}", gifHandler)
		r.Get("/ds/{user_id}.gif", gifHandler)
		r.Get("/ds/{user_id}", gifHandler)
	}

	if cfg.CustomUpload.Enable {
		customGifService := service.NewGIFService(cacheInstance,
			nil,
			petpet.DefaultConfig,
			quantizers.HierarhicalQuantizer{})
		uploadHandler := custom.NewHandler(customGifService, cfg.CustomUpload)

		r.Post("/custom", uploadHandler)
		r.Post("/c", uploadHandler)
	}

	// Set up the server
	var serverAddr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      r,
	}

	if !cfg.EnableHTTP2 {
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}

	if cfg.TLS.Enable && (cfg.TLS.CertFile == "" || cfg.TLS.KeyFile == "") {
		logger.Fatal("TLS is enabled but certFile or keyFile is missing")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	go func() {
		logger.Info("Starting the HTTP server...", zap.String("Address", serverAddr))

		var serveErr error
		switch cfg.TLS.Enable {
		case true:
			serveErr = srv.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		default:
			serveErr = srv.ListenAndServe()
		}

		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Fatal("Server failed:", zap.Error(serveErr))
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
