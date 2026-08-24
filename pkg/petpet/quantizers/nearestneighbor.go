package quantizers

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

const (
	sampleWidth  = 64
	sampleHeight = 4
)

type NearestNeighborQuantizer struct{}

func (NearestNeighborQuantizer) QuantizeImage(img image.Image, numColors int) ([]color.Color, error) {
	smallImg := resize.Resize(sampleWidth, sampleHeight, img, resize.NearestNeighbor)
	bounds := smallImg.Bounds()
	var palette []color.Color

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(palette) >= numColors {
				break
			}
			c := smallImg.At(x, y)
			palette = append(palette, c)
		}
	}

	return palette, nil
}
