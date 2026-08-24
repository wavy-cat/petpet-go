package discord

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wavy-cat/petpet-go/internal/handler/http/utils"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/responses"
	"go.uber.org/zap"
)

func NewHandler(gifService service.GIFService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := middleware.LoggerFromContext(r.Context())
		if !ok {
			http.Error(w, "missing request logger", http.StatusInternalServerError)
			return
		}

		userID, message := verifyDiscordUserID(r)
		if message != "" {
			if userID == "" {
				logger.Warn("Failed to get user ID", zap.String("user_id", userID))
			}
			responses.RespondSoftError(w, message)
			return
		}

		delay, err := utils.ParseDelay(r.URL.Query().Get("delay"))
		if err != nil {
			responses.RespondSoftError(w, "Incorrect delay")
			return
		}

		if r.URL.Query().Get("no-cache") == "true" {
			utils.SetNoCacheHeaders(w)
		}
		gif, err := gifService.GetOrGenerateGif(r.Context(), userID, delay)
		if err != nil {
			logger.Error("Error during GIF generation", zap.Error(err), zap.String("user_id", userID))
			responses.RespondSoftError(w, utils.ParseDiscordError(err))
			return
		}

		if _, err := responses.RespondContent(w, "image/gif", gif); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
	}
}

func verifyDiscordUserID(r *http.Request) (string, string) {
	userID := chi.URLParam(r, "user_id")
	switch {
	case userID == "":
		return userID, "No user ID was specified"
	case strings.EqualFold(userID, "user_id"):
		return userID, "Replace `user_id` in the URL with real Discord user ID 😉"
	default:
		return userID, ""
	}
}
