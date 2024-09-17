package quantizers

import (
	"github.com/Nykakin/quantize"
	"image"
	"image/color"
)

type HierarhicalQuantizer struct{}

func (HierarhicalQuantizer) QuantizeImage(img image.Image, count int) (color.Palette, error) {
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
