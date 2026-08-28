package petpet

import (
	"image"
	"image/color"

	"github.com/Nykakin/quantize"
)

func quantizeImage(img image.Image, count int) (color.Palette, error) {
	quantizer := quantize.NewHierarhicalQuantizer()
	colors, err := quantizer.Quantize(img, count)
	if err != nil {
		return nil, err
	}

	palette := make([]color.Color, len(colors))
	for index, clr := range colors {
		palette[index] = clr
	}

	return palette, nil
}
