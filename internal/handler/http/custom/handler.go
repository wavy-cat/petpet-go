package custom

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Register the JPEG decoder for image.Decode.
	_ "image/png"  // Register the PNG decoder for image.Decode.
	"io"
	"mime/multipart"
	"net/http"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/internal/handler/http/utils"
	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/service"
	"github.com/wavy-cat/petpet-go/pkg/responses"
	"go.uber.org/zap"
	_ "golang.org/x/image/webp" // Register the WebP decoder for image.Decode.
)

const uploadFormName = "image"

var errNoUploadFile = errors.New("no upload file")

func NewHandler(gifService service.GIFService, uploadCfg config.CustomUpload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.LoggerFromContext(r.Context())

		delay, err := utils.ParseDelay(r.URL.Query().Get("delay"))
		if err != nil {
			utils.RespondSoftError(w, "Incorrect delay", logger)
			return
		}

		utils.SetNoCacheHeaders(w)
		img, message := decodeUploadedImage(w, r, uploadCfg, logger)
		if message != "" {
			utils.RespondSoftError(w, message, logger)
			return
		}

		gif, err := gifService.GenerateGifFromImage(r.Context(), img, delay)
		if err != nil {
			logger.Error("Error during GIF generation", zap.Error(err))
			utils.RespondSoftError(w, "Failed to generate GIF", logger)
			return
		}

		if _, err := responses.RespondContent(w, "image/gif", gif); err != nil {
			logger.Error("Error sending response", zap.Error(err))
		}
	}
}

func decodeUploadedImage(
	w http.ResponseWriter,
	r *http.Request,
	uploadCfg config.CustomUpload,
	logger *zap.Logger,
) (image.Image, string) {
	maxUploadSize := boundedUploadSize(uploadCfg.MaxUploadSize)
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	file, err := uploadedFile(r)
	if errors.Is(err, errNoUploadFile) {
		return nil, "No image file was provided"
	}
	if err != nil {
		return nil, fmt.Sprintf(
			"Failed to parse upload. Make sure the file is smaller than %d bytes.",
			uploadCfg.MaxUploadSize,
		)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Error closing upload file", zap.Error(err))
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return nil, fmt.Sprintf(
				"Failed to parse upload. Make sure the file is smaller than %d bytes.",
				uploadCfg.MaxUploadSize,
			)
		}

		return nil, "Unsupported image format. Please upload a PNG, JPEG, or WebP."
	}
	if exceedsPixelLimit(img, uint64(uploadCfg.MaxPixelCount)) {
		return nil, fmt.Sprintf("Image is too large. Maximum allowed size is %d pixels.", uploadCfg.MaxPixelCount)
	}

	return img, ""
}

func uploadedFile(r *http.Request) (*multipart.Part, error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return nil, errNoUploadFile
		}
		if err != nil {
			return nil, err
		}
		if part.FormName() == uploadFormName {
			return part, nil
		}
		if err := part.Close(); err != nil {
			return nil, err
		}
	}
}

func boundedUploadSize(size uint64) int64 {
	const maxInt64 = uint64(1<<63 - 1)
	if size > maxInt64 {
		return int64(maxInt64)
	}

	return int64(size)
}

func exceedsPixelLimit(img image.Image, maxPixelCount uint64) bool {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width <= 0 || height <= 0 {
		return false
	}

	pixelWidth, pixelHeight := uint64(width), uint64(height)
	return pixelWidth > maxPixelCount/pixelHeight
}
