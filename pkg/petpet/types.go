package petpet

import (
	"image"
	"image/color"
)

const (
	DefaultImageSize = 128
	DefaultDelay     = 4
	DefaultDisposal  = 0x02
)

type colorCountedImage struct {
	Image      image.Image
	ColorCount int // Number of colors in the palette
}

type Quantizer interface {
	QuantizeImage(img image.Image, count int) (color.Palette, error)
}
