package answer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondWithPayload отправляет ответ с полезной нагрузкой
func RespondWithPayload(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		respondErr := RespondWithErrorMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		if respondErr != nil {
			return err
		}
		return fmt.Errorf("failed to marshall payload: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err = w.Write(response); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}

// RespondWithErrorMessage отправляет ответ с сообщением об ошибке.
// Обычно используется с кодами 4XX/5XX.
func RespondWithErrorMessage(w http.ResponseWriter, statusCode int, message string) error {
	return RespondWithPayload(w,
		statusCode,
		struct {
			Error string `json:"error"`
		}{
			Error: message,
		})
}

// RespondWithDefaultError отправляет ответ с сообщением об ошибки на основе statusCode.
// Текст берётся из http.StatusText.
func RespondWithDefaultError(w http.ResponseWriter, statusCode int) error {
	return RespondWithErrorMessage(w, statusCode, http.StatusText(statusCode))
}

// RespondOnlyCode отправляет пустой ответ, состоящий только из statusCode.
func RespondOnlyCode(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
