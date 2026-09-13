package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/custom"
	discord_handler "github.com/wavy-cat/petpet-go/internal/handler/http/discord"
	"github.com/wavy-cat/petpet-go/internal/service/animation/gif"
	"github.com/wavy-cat/petpet-go/internal/service/animation/webp"
	"github.com/wavy-cat/petpet-go/pkg/avatarproviders/discord"
	"github.com/wavy-cat/petpet-go/pkg/cache"
)

func addRoutes(r *chi.Mux, cfg config.Config, cacheInstance cache.BytesCache) error {
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("See documentation on GitHub: https://github.com/wavy-cat/petpet-go"))
	})

	if cfg.Discord.Enable {
		if cfg.BotToken == "" {
			return errors.New("discord bot token is required when Discord is enabled")
		}

		// GIF service
		gifHandler := discord_handler.NewHandler(gif.NewGIFService(cacheInstance, discord.NewProvider(cfg.BotToken)))

		r.Get("/discord/{user_id}.gif", gifHandler)
		r.Get("/ds/{user_id}.gif", gifHandler)

		// WebP service
		webpHandler := discord_handler.NewHandler(webp.NewWebPService(cacheInstance, discord.NewProvider(cfg.BotToken)))

		r.Get("/discord/{user_id}.webp", webpHandler)
		r.Get("/discord/{user_id}", webpHandler)
		r.Get("/ds/{user_id}.webp", webpHandler)
		r.Get("/ds/{user_id}", webpHandler)
	}

	if cfg.CustomUpload.Enable {
		// GIF Service
		gifHandler := custom.NewHandler(gif.NewGIFService(cacheInstance, nil), cfg.CustomUpload)

		r.Post("/custom/gif", gifHandler)
		r.Post("/c/gif", gifHandler)

		// WebP service
		webpHandler := custom.NewHandler(webp.NewWebPService(cacheInstance, nil), cfg.CustomUpload)

		r.Post("/custom/webp", webpHandler)
		r.Post("/custom", webpHandler)
		r.Post("/c/webp", webpHandler)
		r.Post("/c", webpHandler)
	}

	return nil
}
