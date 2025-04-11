package responses

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondContent(t *testing.T) {
	w := httptest.NewRecorder()
	content := []byte("test content")
	contentType := "text/plain"

	_, err := RespondContent(w, contentType, content)
	if err != nil {
		t.Errorf("RespondContent returned error: %v", err)
	}

	if w.Header().Get("Content-Type") != contentType {
		t.Errorf("unexpected Content-Type header: got %s, want %s", w.Header().Get("Content-Type"), contentType)
	}

	if !bytes.Equal(w.Body.Bytes(), content) {
		t.Errorf("unexpected response body: got %v, want %v", w.Body.Bytes(), content)
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()

	err := RespondError(w)
	if err != nil {
		t.Errorf("RespondError returned error: %v", err)
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusInternalServerError)
	}

	headers := map[string]string{
		"Cache-Control": "no-store, no-cache, must-revalidate, private",
		"Pragma":        "no-cache",
		"Expires":       "0",
	}

	for k, v := range headers {
		if w.Header().Get(k) != v {
			t.Errorf("%s header = %s, want %s", k, w.Header().Get(k), v)
		}
	}
}

func TestRespondSoftError(t *testing.T) {
	w := httptest.NewRecorder()
	const details = "test error details"

	err := RespondSoftError(w, details)
	if err != nil {
		t.Errorf("RespondSoftError returned error: %v", err)
	}

	headers := map[string]string{
		"Cache-Control": "no-store, no-cache, must-revalidate, private",
		"Pragma":        "no-cache",
		"Expires":       "0",
		"Content-Type":  "text/html",
	}

	for k, v := range headers {
		if w.Header().Get(k) != v {
			t.Errorf("%s header = %s, want %s", k, w.Header().Get(k), v)
		}
	}

	if !strings.Contains(w.Body.String(), details) {
		t.Errorf("Response body does not contain details: %s", details)
	}
}
