package petpet

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"sync"
)

func pasteImageRGBA(dest *image.RGBA, src image.Image, offsetX, offsetY int) {
	draw.Draw(dest, src.Bounds().Add(image.Pt(offsetX, offsetY)), src, image.Point{}, draw.Over)
}

func MakeAPNG(source io.Reader, writer io.Writer, config Config) error {
	baseImg, err := png.Decode(source)
	if err != nil {
		return err
	}

	const frames = 10
	baseImg = resizeImage(baseImg, config.Width, config.Height)
	images := make([]*image.RGBA, frames)

	var wg sync.WaitGroup
	wg.Add(frames)

	for i := 0; i < frames; i++ {
		go func(i int) {
			defer wg.Done()
			canvas := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

			squeeze := float64(i)
			if i >= frames/2 {
				squeeze = float64(frames - i)
			}

			var (
				scaleX  = 0.8 + squeeze*0.02
				scaleY  = 0.8 - squeeze*0.05
				offsetX = int((1 - scaleX) * float64(config.Width) * 0.5)
				offsetY = int((1 - scaleY) * float64(config.Height))
			)

			resizedImg := resizeImage(baseImg, int(float64(config.Width)*scaleX), int(float64(config.Height)*scaleY))
			pasteImageRGBA(canvas, resizedImg, offsetX, offsetY)

			petFrame := resizeImage(hands[i], config.Width, config.Height)
			pasteImageRGBA(canvas, petFrame, 0, 0)

			images[i] = canvas
		}(i)
	}

	wg.Wait()

	return exportAPNG(writer, images, uint16(config.Delay))
}
