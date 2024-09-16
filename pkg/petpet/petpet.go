package petpet

import (
	"bytes"
	"github.com/nfnt/resize"
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

// Генерирует палитру на основе цветов изображения.
// Первый цвет всегда прозрачный (0, 0, 0, 0).
func generatePalette(img image.Image, numColors int) []color.Color {
	smallImg := resize.Resize(64, 0, img, resize.NearestNeighbor)
	bounds := smallImg.Bounds()
	colorMap := make(map[color.Color]struct{})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := smallImg.At(x, y)
			colorMap[c] = struct{}{}
		}
	}

	var palette []color.Color
	for c := range colorMap {
		palette = append(palette, c)
		if len(palette) >= numColors {
			break
		}
	}

	palette[0] = color.RGBA{}

	return palette
}

func createTransparentImage(baseImg image.Image, width, height int) *image.Paletted {
	rect := image.Rect(0, 0, width, height)
	return image.NewPaletted(rect, generatePalette(baseImg, 256))
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

	for i := 0; i < frames; i++ {
		canvas := createTransparentImage(baseImg, width, height)

		squeeze := float64(i)
		if i >= frames/2 {
			squeeze = float64(frames - i)
		}

		scaleX := 0.8 + squeeze*0.02
		scaleY := 0.8 - squeeze*0.05
		offsetX := int((1 - scaleX) * float64(width) * 0.5)
		offsetY := int((1 - scaleY) * float64(height))

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
