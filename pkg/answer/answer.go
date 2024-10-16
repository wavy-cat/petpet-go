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

// RespondHTMLError отправляет ошибку в виде HTML с meta-тегами
func RespondHTMLError(w http.ResponseWriter, title, details string) (int, error) {
	const body = `<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>PetPet</title>
			<meta content="%s" property="og:title"/>
			<meta content="%s" property="og:description"/>
			<meta content="https://github.com/wavy-cat/petpet-go" property="og:url"/>
			<meta content="#EE204D" data-react-helmet="true" name="theme-color"/>
			<style>
			body {
				color: white;
				background-color: black;
			}
			</style>
		</head>
		<body>
			<p>%s</p>
		</body>
		</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	responseBody := fmt.Sprintf(body, title, details, details)
	return w.Write([]byte(responseBody))
}
