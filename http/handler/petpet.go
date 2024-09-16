package handler

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/internal/answer"
	"github.com/wavy-cat/petpet-go/pkg/avatar"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
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
		err := answer.RespondWithErrorMessage(w, http.StatusBadRequest, "User ID not sent")
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Получаем аватар по ID
	avatarImage, err := avatar.GetAvatarFromID(userId)
	if err != nil {
		logger.Warn("Failed to get user avatar",
			zap.Error(err), zap.String("User ID", userId))
		err := answer.RespondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		if err != nil {
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
			err := answer.RespondWithErrorMessage(w, http.StatusBadRequest, err.Error())
			if err != nil {
				logger.Error("Error sending response", zap.Error(err))
			}
			return
		}
	}

	// Генерируем гифку
	config := petpet.DefaultConfig
	config.Delay = delay
	avatarReader := bytes.NewReader(avatarImage)

	gif, err := petpet.MakeGif(avatarReader, config)
	if err != nil {
		logger.Error("Failed to generate gif",
			zap.Error(err), zap.String("User ID", userId))
		err := answer.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Отправляем гифку
	buf := make([]byte, 1024)
	for {
		n, err := gif.Read(buf)
		if err != nil && err != io.EOF {
			logger.Error("Error reading gif",
				zap.Error(err), zap.String("User ID", userId))
			err := answer.RespondWithDefaultError(w, http.StatusInternalServerError)
			if err != nil {
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
