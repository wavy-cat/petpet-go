package petpet

import (
	"image"
	"image/color"
)

const (
	defaultImageSize = 128
	defaultDelay     = 4
	defaultDisposal  = 0x02
)

type Config struct {
	Width    int  // Recommend 128
	Height   int  // Recommend 128
	Delay    int  // Recommend 2-10
	Disposal byte // Recommend 0x02
}

var DefaultConfig = Config{
	Width:    defaultImageSize,
	Height:   defaultImageSize,
	Delay:    defaultDelay,
	Disposal: defaultDisposal,
}

type colorCountedImage struct {
	Image      image.Image
	ColorCount int // Number of colors in the palette
}

type Quantizer interface {
	QuantizeImage(img image.Image, count int) (color.Palette, error)
}
