package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	httpServer *http.Server
	tlsConfig  config.ServerTLS
}

func NewServer(cfg config.Config, logger *zap.Logger, cacheInstance cache.BytesCache) (Server, error) {
	r := chi.NewRouter()

	addMiddlewares(r, cfg, logger)

	if err := addRoutes(r, cfg, cacheInstance); err != nil {
		return nil, fmt.Errorf("error when adding routes: %w", err)
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

	if cfg.TLS.Enable && (cfg.TLS.CertFile == "" || cfg.TLS.KeyFile == "") {
		return nil, errors.New("tls is enabled but certFile or keyFile is missing")
	}

	return &server{
		httpServer: srv,
		tlsConfig:  cfg.TLS,
	}, nil
}

func (s *server) Start() error {
	if s.tlsConfig.Enable {
		return s.httpServer.ListenAndServeTLS(s.tlsConfig.CertFile, s.tlsConfig.KeyFile)
	}
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
