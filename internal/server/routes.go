package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/custom"
	discord_handler "github.com/wavy-cat/petpet-go/internal/handler/http/discord"
	"github.com/wavy-cat/petpet-go/internal/service"
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

		discordGifService := service.NewGIFService(cacheInstance,
			discord.NewProvider(cfg.BotToken))
		gifHandler := discord_handler.NewHandler(discordGifService)

		r.Get("/discord/{user_id}.gif", gifHandler)
		r.Get("/discord/{user_id}", gifHandler)
		r.Get("/ds/{user_id}.gif", gifHandler)
		r.Get("/ds/{user_id}", gifHandler)
	}

	if cfg.CustomUpload.Enable {
		customGifService := service.NewGIFService(cacheInstance,
			nil)
		uploadHandler := custom.NewHandler(customGifService, cfg.CustomUpload)

		r.Post("/custom", uploadHandler)
		r.Post("/c", uploadHandler)
	}

	return nil
}
