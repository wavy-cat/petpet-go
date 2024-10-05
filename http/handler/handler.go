package handler

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/internal/answer"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/discord"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{}

func (Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Получаем логгер
	logger := r.Context().Value("logger").(*zap.Logger)

	// Получаем ID пользователя
	userId, ok := mux.Vars(r)["user_id"]
	if !ok {
		logger.Warn("Failed to get user ID")
		if err := answer.RespondWithErrorMessage(w, http.StatusBadRequest, "User ID not sent"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	if strings.ToLower(userId) == "user_id" {
		_, err := answer.RespondHTMLError(w, "Misuse",
			"Replace user_id in the URL with real Discord user ID 😉")
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Получаем delay из параметров
	delayParam := strings.TrimSpace(r.URL.Query().Get("delay"))
	var delay int

	switch strings.TrimSpace(r.URL.Query().Get("delay")) {
	case "":
		delay = 2
	default:
		var err error
		delay, err = strconv.Atoi(delayParam)
		if err != nil {
			if _, err := answer.RespondHTMLError(w, "Incorrect delay", err.Error()); err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}
	}

	cacheObj, cacheOk := r.Context().Value("cache").(cache.BytesCache) // Getting the cache object
	bot, ok := r.Context().Value("bot").(*discord.Bot)                 // Getting the bot object

	// Получаем аватар по ID
	var (
		avatarImage []byte
		avatarId    string
	)

	switch ok {
	case true:
		user, err := bot.NewUserById(userId)
		if err != nil {
			logger.Warn("Failed to get user", zap.Error(err), zap.String("User ID", userId))
			title, status := checkError(err)

			_, err = answer.RespondHTMLError(w, title, status)
			if err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}

		// Checking whether the user has an avatar
		if user.Avatar == nil {
			_, err = answer.RespondHTMLError(w, "Not found", "Avatar not found")
			if err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}

		// Getting the finished image from the cache
		if cacheOk {
			ok, err := responseFromCache(w, cacheObj, *user.Avatar)
			if err != nil {
				logger.Error("Error sending response from cache", zap.Error(err))
			}
			if ok {
				return
			}
		}

		// Getting the user's avatar
		avatarImage, err = user.GetAvatar()
		if err != nil {
			logger.Warn("Failed to load user avatar", zap.Error(err), zap.String("User ID", userId))
			_, err = answer.RespondHTMLError(w, "Unknown Error", "Something went wrong")
			if err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}

		avatarId = *user.Avatar
	case false:
		var err error
		avatarImage, avatarId, err = getAvatarUsingCdev(userId)
		if err != nil {
			logger.Warn("Failed to load user avatar", zap.Error(err), zap.String("User ID", userId))
			_, err = answer.RespondHTMLError(w, "Error getting avatar", err.Error())
			if err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}

		// Getting the finished image from the cache
		if cacheOk {
			ok, err := responseFromCache(w, cacheObj, avatarId)
			if err != nil {
				logger.Error("Error sending response from cache", zap.Error(err))
			}
			if ok {
				return
			}
		}
	}

	// Получаем no-cache из параметров
	switch r.URL.Query().Get("no-cache") {
	case "true":
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache") // Для совместимости со старыми браузерами
		w.Header().Set("Expires", "0")       // Для совместимости со старыми браузерами
	default:
		w.Header().Set("Cache-Control", "max-age=900")
	}

	// Генерируем гифку
	config := petpet.DefaultConfig
	config.Delay = delay
	avatarReader := bytes.NewReader(avatarImage)

	gif, err := petpet.MakeGif(avatarReader, config, quantizers.HierarhicalQuantizer{})
	if err != nil {
		logger.Error("Failed to generate gif", zap.Error(err), zap.String("User ID", userId))
		if err := answer.RespondWithDefaultError(w, http.StatusInternalServerError); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Устанавливаем Content-Type
	w.Header().Set("Content-Type", "image/gif")

	// Отправляем гифку
	data, err := io.ReadAll(gif)
	if err != nil {
		logger.Error("Failed to read GIF Reader", zap.Error(err))
		_, err = answer.RespondHTMLError(w, "Internal Server Error", "Something went wrong")
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	_, err = w.Write(data)
	if err != nil {
		logger.Error("Error sending response", zap.Error(err))
	}

	// Adding an image to the cache
	if cacheOk {
		err = cacheObj.Push(avatarId, data)
		if err != nil {
			logger.Error("Failed to write data to cache", zap.Error(err))
		}
	}
}
