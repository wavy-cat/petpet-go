package petpet

import (
	"image"

	"github.com/nfnt/resize"
)

func resizeImage(img image.Image, newWidth, newHeight int) image.Image {
	return resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
}
