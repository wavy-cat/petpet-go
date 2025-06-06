package petpet

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"sync"
)

// Создаёт палитру цветов на основе переданных изображений.
// Сумма цветов всех изображений не должна превышать 256 с `addTransparent` в значении false
// или `255` если установлено `true`.
func createPalette(addTransparent bool, quantizer Quantizer, images ...colorCountedImage) (color.Palette, error) {
	palette := make([]color.Color, 0, 256)
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

	if len(palette) > 256 {
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

// MakeGif генерирует pet-pet гифку.
// `source` должен быть типом io.Reader и содержать PNG изображение.
func MakeGif(source io.Reader, w io.Writer, config Config, quantizer Quantizer) error {
	var (
		width    = config.Width
		height   = config.Height
		delay    = config.Delay
		disposal = config.Disposal
	)
	const frames = 10

	baseImg, err := png.Decode(source)
	if err != nil {
		return err
	}

	if size := baseImg.Bounds().Size(); size.X != width || size.Y != width {
		baseImg = resizeImage(baseImg, width, height)
	}

	var images = make([]*image.Paletted, frames)

	basePalette, err := createPalette(
		true,
		quantizer,
		[]colorCountedImage{
			{
				Image:      baseImg,
				ColorCount: 240,
			},
			{
				Image:      hands[0],
				ColorCount: 15,
			}}...)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(frames)

	for i := range frames {
		go func(i int) {
			squeeze := float64(i)
			if i >= frames/2 {
				squeeze = float64(frames - i)
			}

			var (
				scaleX  = 0.8 + squeeze*0.02
				scaleY  = 0.8 - squeeze*0.05
				offsetX = int(((1-scaleX)*0.5 + 0.1) * float64(width))
				offsetY = int(((1 - scaleY) - 0.08) * float64(height))
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
		delays    = make([]int, frames)
		disposals = make([]byte, frames)
	)

	for i := range frames {
		delays[i] = delay
		disposals[i] = disposal
	}

	wg.Wait()

	return exportGIF(w, images, delays, disposals)
}
