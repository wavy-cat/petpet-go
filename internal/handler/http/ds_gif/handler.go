package ds_gif

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/petpet-go/internal/handler/http/utils"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/answer"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	gifService service.GIFService
}

func NewHandler(gifService service.GIFService) *Handler {
	return &Handler{gifService: gifService}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*zap.Logger)

	// Getting the user ID
	userId, ok := mux.Vars(r)["user_id"]
	if !ok {
		logger.Warn("Failed to get user ID", zap.String("user_id", userId))
		if err := answer.RespondWithErrorMessage(w, http.StatusBadRequest, "User ID not sent"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Getting delay
	delay, err := utils.ParseDelay(r.URL.Query().Get("delay"))
	if err != nil {
		if _, err := answer.RespondHTMLError(w, "Incorrect delay", err.Error()); err != nil {
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
	gif, err := h.gifService.GetOrGenerateGif(ctx, userId, "discord", delay)
	if err != nil {
		logger.Warn("Error during GIF generation", zap.Error(err))
		title, description := utils.ParseError(err)
		if _, err := answer.RespondHTMLError(w, title, description); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	// Setting Content-Type
	w.Header().Set("Content-Type", "image/gif")

	// Returning the result
	_, err = w.Write(gif)
	if err != nil {
		logger.Error("Error sending response", zap.Error(err))
	}
}
