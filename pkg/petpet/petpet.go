package petpet

import (
	"bytes"
	"errors"
	"github.com/nfnt/resize"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
)

func exportGIF(images []*image.Paletted, delays []int, disposals []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	g := &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: 0,
		Disposal:  disposals,
	}

	err := gif.EncodeAll(&buf, g)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func resizeImage(img image.Image, newWidth, newHeight int) image.Image {
	return resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
}

// Создаёт палитру цветов на основе переданных изображений.
// Сумма цветов всех изображений не должна превышать 256 с `addTransparent` в значении false
// или `255` если установлено `true`.
func createPalette(addTransparent bool, quantizer quantizer, images ...colorCountedImage) (color.Palette, error) {
	palette := make([]color.Color, 0, 256)
	if addTransparent {
		palette = []color.Color{color.RGBA{}}
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
func MakeGif(source io.Reader, config Config) (io.Reader, error) {
	baseImg, err := png.Decode(source)
	if err != nil {
		return nil, err
	}

	var (
		width    = config.Width
		height   = config.Height
		delay    = config.Delay
		disposal = config.Disposal
	)
	const frames = 10

	baseImg = resizeImage(baseImg, width, height)
	var (
		images    []*image.Paletted
		delays    []int
		disposals []byte
	)

	basePalette, err := createPalette(
		true,
		quantizers.HierarhicalQuantizer{},
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
		return nil, err
	}

	for i := 0; i < frames; i++ {
		canvas := createTransparentImage(width, height, basePalette)

		squeeze := float64(i)
		if i >= frames/2 {
			squeeze = float64(frames - i)
		}

		var (
			scaleX  = 0.8 + squeeze*0.02
			scaleY  = 0.8 - squeeze*0.05
			offsetX = int((1 - scaleX) * float64(width) * 0.5)
			offsetY = int((1 - scaleY) * float64(height))
		)

		resizedImg := resizeImage(baseImg, int(float64(width)*scaleX), int(float64(height)*scaleY))
		pasteImage(canvas, resizedImg, offsetX, offsetY)

		petFrame := resizeImage(hands[i], width, height)
		pasteImage(canvas, petFrame, 0, 0)

		images = append(images, canvas)
		delays = append(delays, delay)
		disposals = append(disposals, disposal)
	}

	data, err := exportGIF(images, delays, disposals)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
