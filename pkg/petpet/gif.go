package petpet

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io"
	"sync"
)

const (
	maxPaletteColors      = 256
	baseImageColorCount   = 240
	handImageColorCount   = 15
	horizontalImageOffset = 0.1
	verticalImageOffset   = 0.08
	animationFrameCount   = 10
)

func createPalette(addTransparent bool, quantizer Quantizer, images ...colorCountedImage) (color.Palette, error) {
	palette := make([]color.Color, 0, maxPaletteColors)
	if addTransparent {
		palette = append(palette, color.RGBA{})
	}

	for _, val := range images {
		imgPalette, err := quantizer.QuantizeImage(val.Image, val.ColorCount)
		if err != nil {
			return nil, err
		}
		palette = append(palette, imgPalette...)
	}

	if len(palette) > maxPaletteColors {
		return nil, errors.New("the palette has more than 256 colors")
	}

	return palette, nil
}

func createTransparentImage(width, height int, palette color.Palette) *image.Paletted {
	rect := image.Rect(0, 0, width, height)
	return image.NewPaletted(rect, palette)
}

func pasteImage(dest *image.Paletted, src image.Image, offsetX, offsetY int) {
	draw.Draw(dest, src.Bounds().Add(image.Pt(offsetX, offsetY)), src, image.Point{}, draw.Over)
}

// MakeGif generates a pet-pet gif.
func MakeGif(baseImg image.Image, w io.Writer, config Config, quantizer Quantizer) error {
	var (
		width    = config.Width
		height   = config.Height
		delay    = config.Delay
		disposal = config.Disposal
	)
	if size := baseImg.Bounds().Size(); size.X != width || size.Y != width {
		baseImg = resizeImage(baseImg, width, height)
	}

	var images = make([]*image.Paletted, animationFrameCount)

	basePalette, err := createPalette(
		true,
		quantizer,
		[]colorCountedImage{
			{
				Image:      baseImg,
				ColorCount: baseImageColorCount,
			},
			{
				Image:      hands[0],
				ColorCount: handImageColorCount,
			}}...)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(animationFrameCount)

	for i := range animationFrameCount {
		go func(i int) {
			squeeze := float64(i)
			if i >= animationFrameCount/2 {
				squeeze = float64(animationFrameCount - i)
			}

			var (
				scaleX  = 0.8 + squeeze*0.02
				scaleY  = 0.8 - squeeze*0.05
				offsetX = int(((1-scaleX)*0.5 + horizontalImageOffset) * float64(width))
				offsetY = int(((1 - scaleY) - verticalImageOffset) * float64(height))
			)

			canvas := createTransparentImage(width, height, basePalette)

			resizedImg := resizeImage(baseImg, int(float64(width)*scaleX), int(float64(height)*scaleY))
			pasteImage(canvas, resizedImg, offsetX, offsetY)

			petFrame := resizeImage(hands[i], width, height)
			pasteImage(canvas, petFrame, 0, 0)

			images[i] = canvas

			wg.Done()
		}(i)
	}

	var (
		delays    = make([]int, animationFrameCount)
		disposals = make([]byte, animationFrameCount)
	)

	for i := range animationFrameCount {
		delays[i] = delay
		disposals[i] = disposal
	}

	wg.Wait()

	return exportGIF(w, images, delays, disposals)
}
