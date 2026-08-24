package server

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"go.uber.org/zap"
)

func addMiddlewares(r *chi.Mux, cfg config.Config, logger *zap.Logger) {
	r.Use(middleware.Logger(logger))

	if cfg.Heartbeat.Enable {
		r.Use(chiMiddleware.Heartbeat(cfg.Heartbeat.Path))
	}

	if cfg.Throttle.Enable {
		r.Use(chiMiddleware.ThrottleBacklog(
			cfg.Throttle.Limit,
			cfg.Throttle.Backlog,
			cfg.Throttle.BacklogTimeout))
	}
}
