package responses

import (
	"fmt"
	"net/http"
	"strings"
)

// RespondContent sends any content as well as the type of content in the header.
func RespondContent(w http.ResponseWriter, contentType string, content []byte) (int, error) {
	w.Header().Set("Content-Type", contentType)

	return w.Write(content)
}

// RespondError sends error 500 with no details.
func RespondError(w http.ResponseWriter) error {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

// RespondSoftError sends the error as HTML with meta tags, which Discord displays as an embed.
// Disables cache in headers.
func RespondSoftError(w http.ResponseWriter, details string) error {
	const body = `<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>PetPet - Error</title>
			<meta content="Error" property="og:title"/>
			<meta content="%s" property="og:description"/>
			<meta content="https://github.com/wavy-cat/petpet-go" property="og:url"/>
			<meta content="#EE204D" data-react-helmet="true" name="theme-color"/>
		</head>
		<body>
			<p>%s</p>
		</body>
		</html>`

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	responseBody := fmt.Sprintf(body, details, details)
	responseBody = strings.ReplaceAll(responseBody, "\t", "")
	responseBody = strings.ReplaceAll(responseBody, "\n", "")

	_, err := RespondContent(w, "text/html", []byte(responseBody))
	return err
}
