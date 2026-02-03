package custom

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/utils"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/responses"
	"go.uber.org/zap"
)

const uploadFormName = "image"

type Handler struct {
	gifService service.GIFService
	uploadCfg  config.CustomUpload
}

func NewHandler(gifService service.GIFService, uploadCfg config.CustomUpload) *Handler {
	return &Handler{
		gifService: gifService,
		uploadCfg:  uploadCfg,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(middleware.LoggerKey).(*zap.Logger)

	delay, err := utils.ParseDelay(r.URL.Query().Get("delay"))
	if err != nil {
		if err := responses.RespondSoftError(w, "Incorrect delay"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	maxUploadSize := int64(h.uploadCfg.MaxUploadSize)
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		if err := responses.RespondSoftError(w, fmt.Sprintf("Failed to parse upload. Make sure the file is smaller than %d bytes.", h.uploadCfg.MaxUploadSize)); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	file, _, err := r.FormFile(uploadFormName)
	if err != nil {
		if err := responses.RespondSoftError(w, "No image file was provided"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing upload file", zap.Error(err))
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		if err := responses.RespondSoftError(w, "Unsupported image format. Please upload a PNG or JPEG."); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	bounds := img.Bounds()
	pixelCount := int64(bounds.Dx()) * int64(bounds.Dy())
	if pixelCount > int64(h.uploadCfg.MaxPixelCount) {
		if err := responses.RespondSoftError(w, fmt.Sprintf("Image is too large. Maximum allowed size is %d pixels.", h.uploadCfg.MaxPixelCount)); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	gif, err := h.gifService.GenerateGifFromImage(r.Context(), img, delay)
	if err != nil {
		logger.Error("Error during GIF generation", zap.Error(err))
		if err := responses.RespondSoftError(w, "Failed to generate GIF"); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
		return
	}

	_, err = responses.RespondContent(w, "image/gif", gif)
	if err != nil {
		logger.Error("Error sending response", zap.Error(err))
	}
}
