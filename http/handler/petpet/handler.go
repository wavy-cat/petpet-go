package petpet

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/internal/answer"
	"github.com/wavy-cat/petpet-go/pkg/avatar"
	"github.com/wavy-cat/petpet-go/pkg/discord"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
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

	// Получаем объект бота
	bot, ok := r.Context().Value("bot").(*discord.Bot)

	// Получаем аватар по ID
	var avatarImage []byte
	var err error

	switch ok {
	case true:
		avatarImage, err = GetAvatarUsingBot(bot, userId)
	case false:
		avatarImage, err = avatar.GetAvatarFromID(userId)
	}

	if err != nil {
		logger.Warn("Failed to get user", zap.Error(err), zap.String("User ID", userId))
		if err := answer.RespondWithErrorMessage(w, http.StatusInternalServerError, err.Error()); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Получаем delay из параметров
	delayParam := r.URL.Query().Get("delay")
	var delay int
	switch delayParam {
	case "":
		delay = 2
	default:
		delay, err = strconv.Atoi(delayParam)
		if err != nil {
			if err := answer.RespondWithErrorMessage(w, http.StatusBadRequest, err.Error()); err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}
	}

	// Получаем no-cache из параметров
	switch r.URL.Query().Get("no-cache") {
	case "true":
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache") // Для совместимости со старыми браузерами
		w.Header().Set("Expires", "0")       // Для совместимости со старыми браузерами
	default:
		w.Header().Set("Cache-Control", "max-age=300")
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
	buf := make([]byte, 1024)
	for {
		n, err := gif.Read(buf)
		if err != nil && err != io.EOF {
			logger.Error("Error reading gif", zap.Error(err), zap.String("User ID", userId))
			if err := answer.RespondWithDefaultError(w, http.StatusInternalServerError); err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}

		if err == io.EOF || n == 0 {
			break
		}

		if _, err := w.Write(buf); err != nil {
			logger.Error("Error sending response", zap.Error(err))
			break
		}
	}
}
