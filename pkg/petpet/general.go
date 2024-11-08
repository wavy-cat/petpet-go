package petpet

import (
	"github.com/nfnt/resize"
	"image"
)

func resizeImage(img image.Image, newWidth, newHeight int) image.Image {
	return resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
}
