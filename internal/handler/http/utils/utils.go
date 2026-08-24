package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/wavy-cat/petpet-go/pkg/responses"
	"go.uber.org/zap"
)

const defaultDelay = 3

func ParseDelay(value string) (int, error) {
	value = strings.TrimSpace(value)
	switch value {
	case "":
		return defaultDelay, nil
	default:
		return strconv.Atoi(value)
	}
}

func ParseDiscordError(err error) string {
	switch {
	case strings.Contains(err.Error(), "10013"):
		return "User not found"
	case strings.Contains(err.Error(), "50035"):
		return "Incorrect user ID. Check ID for correctness"
	case strings.Contains(err.Error(), "get avatar error: not exists"):
		return "User has no avatar"
	}

	return "Something went wrong"
}

func SetNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func RespondSoftError(w http.ResponseWriter, details string, logger *zap.Logger) {
	if err := responses.RespondSoftError(w, details); err != nil {
		logger.Error("Error sending response", zap.Error(err))
	}
}
