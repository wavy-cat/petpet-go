package ds_apng

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wavy-cat/petpet-go/internal/handler/http/utils"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/responses"
	"go.uber.org/zap"
)

type Handler struct {
	apngService service.APNGService
	transport   *http.Transport
}

func NewHandler(apngService service.APNGService, transport *http.Transport) *Handler {
	return &Handler{
		apngService: apngService,
		transport:   transport,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*zap.Logger)

	// Getting the user ID
	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		logger.Warn("Failed to get user ID", zap.String("user_id", userId))
		if err := responses.RespondSoftError(w, "No user ID was specified"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	if strings.ToLower(userId) == "user_id" {
		if err := responses.RespondSoftError(w, "Replace `user_id` in the URL with real Discord user ID 😉"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Getting delay
	delay, err := utils.ParseDelay(r.URL.Query().Get("delay"))
	if err != nil {
		if err := responses.RespondSoftError(w, "Incorrect delay"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Setting caching policies
	switch r.URL.Query().Get("no-cache") {
	case "true":
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache") // For compatibility with older browsers
		w.Header().Set("Expires", "0")       // For compatibility with older browsers
	default:
		w.Header().Set("Cache-Control", "max-age=900")
	}

	// Calling the service to generate GIF
	ctx := context.WithValue(context.Background(), "logger", logger)
	ctx = context.WithValue(ctx, "transport", h.transport)
	gif, err := h.apngService.GetOrGenerateAPNG(ctx, userId, "discord", delay)
	if err != nil {
		logger.Warn("Error during APNG generation", zap.Error(err))

		errDetails := utils.ParseDiscordError(err)
		if err := responses.RespondSoftError(w, errDetails); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Returning the result
	_, err = responses.RespondContent(w, "image/gif", gif)
	if err != nil {
		logger.Error("Error sending response", zap.Error(err))
	}
}
