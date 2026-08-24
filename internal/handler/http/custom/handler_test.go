package custom_test

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/custom"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"go.uber.org/zap"
)

const uploadFormName = "image"

type fakeGIFService struct {
	called          bool
	lastDelay       int
	response        []byte
	responseErr     error
	lastImageWidth  int
	lastImageHeight int
}

func (f *fakeGIFService) GetOrGenerateGif(_ context.Context, _ string, _ int) ([]byte, error) {
	return nil, nil
}

func (f *fakeGIFService) GenerateGifFromImage(_ context.Context, img image.Image, delay int) ([]byte, error) {
	f.called = true
	f.lastDelay = delay
	bounds := img.Bounds()
	f.lastImageWidth = bounds.Dx()
	f.lastImageHeight = bounds.Dy()
	return f.response, f.responseErr
}

func newRequestWithLogger(t *testing.T, req *http.Request) *http.Request {
	t.Helper()
	ctx := context.WithValue(req.Context(), middleware.LoggerKey, zap.NewNop())
	return req.WithContext(ctx)
}

func buildMultipartRequest(t *testing.T, fieldName, fileName string, data []byte) (*http.Request, error) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if fieldName != "" {
		part, err := writer.CreateFormFile(fieldName, fileName)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(data); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/custom", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestHandlerSuccess(t *testing.T) {
	t.Parallel()

	svc := &fakeGIFService{response: []byte("gifdata")}
	handler := custom.NewHandler(svc, config.CustomUpload{MaxUploadSize: 1024 * 1024, MaxPixelCount: 100})

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	req, err := buildMultipartRequest(t, uploadFormName, "image.png", imgBuf.Bytes())
	if err != nil {
		t.Fatalf("build multipart request: %v", err)
	}
	req = newRequestWithLogger(t, req)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if !svc.called {
		t.Fatal("expected gif service to be called")
	}
	if got := recorder.Header().Get("Content-Type"); got != "image/gif" {
		t.Fatalf("expected image/gif content type, got %q", got)
	}
	if body := recorder.Body.String(); body != "gifdata" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func TestHandlerInvalidDelay(t *testing.T) {
	t.Parallel()

	handler := custom.NewHandler(&fakeGIFService{}, config.CustomUpload{MaxUploadSize: 1024, MaxPixelCount: 10})
	req, err := buildMultipartRequest(t, uploadFormName, "image.png", []byte("data"))
	if err != nil {
		t.Fatalf("build multipart request: %v", err)
	}
	req = newRequestWithLogger(t, req)
	req.URL.RawQuery = "delay=abc"

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if body := recorder.Body.String(); !strings.Contains(body, "Incorrect delay") {
		t.Fatalf("expected delay error, got %q", body)
	}
}

func TestHandlerMissingFile(t *testing.T) {
	t.Parallel()

	handler := custom.NewHandler(&fakeGIFService{}, config.CustomUpload{MaxUploadSize: 1024, MaxPixelCount: 10})
	req, err := buildMultipartRequest(t, "", "", nil)
	if err != nil {
		t.Fatalf("build multipart request: %v", err)
	}
	req = newRequestWithLogger(t, req)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if body := recorder.Body.String(); !strings.Contains(body, "No image file was provided") {
		t.Fatalf("expected missing file error, got %q", body)
	}
}

func TestHandlerUnsupportedFormat(t *testing.T) {
	t.Parallel()

	handler := custom.NewHandler(&fakeGIFService{}, config.CustomUpload{MaxUploadSize: 1024, MaxPixelCount: 10})
	req, err := buildMultipartRequest(t, uploadFormName, "image.txt", []byte("not an image"))
	if err != nil {
		t.Fatalf("build multipart request: %v", err)
	}
	req = newRequestWithLogger(t, req)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if body := recorder.Body.String(); !strings.Contains(body, "Unsupported image format") {
		t.Fatalf("expected unsupported format error, got %q", body)
	}
}

func TestHandlerPixelLimit(t *testing.T) {
	t.Parallel()

	handler := custom.NewHandler(&fakeGIFService{}, config.CustomUpload{MaxUploadSize: 1024 * 1024, MaxPixelCount: 3})
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	req, err := buildMultipartRequest(t, uploadFormName, "image.png", imgBuf.Bytes())
	if err != nil {
		t.Fatalf("build multipart request: %v", err)
	}
	req = newRequestWithLogger(t, req)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if body := recorder.Body.String(); !strings.Contains(body, "Maximum allowed size is 3 pixels") {
		t.Fatalf("expected pixel limit error, got %q", body)
	}
}

func TestHandlerUploadTooLarge(t *testing.T) {
	t.Parallel()

	handler := custom.NewHandler(&fakeGIFService{}, config.CustomUpload{MaxUploadSize: 10, MaxPixelCount: 10})
	req, err := buildMultipartRequest(t, uploadFormName, "image.png", []byte("0123456789abcdef"))
	if err != nil {
		t.Fatalf("build multipart request: %v", err)
	}
	req = newRequestWithLogger(t, req)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if body := recorder.Body.String(); !strings.Contains(body, "Failed to parse upload") {
		t.Fatalf("expected upload size error, got %q", body)
	}
}
