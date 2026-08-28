package petpet

import (
	"image"
	"image/draw"
	"sync"
)

const (
	horizontalImageOffset = 0.1
	verticalImageOffset   = 0.08
	animationFrameCount   = 10
)

func prepareAnimationBase(baseImg image.Image, config Config) image.Image {
	if size := baseImg.Bounds().Size(); size.X != config.Width || size.Y != config.Width {
		return resizeImage(baseImg, config.Width, config.Height)
	}

	return baseImg
}

func createAnimationFrames(
	baseImg image.Image,
	config Config,
	createCanvas func(width, height int) draw.Image,
) []image.Image {
	width := config.Width
	height := config.Height
	images := make([]image.Image, animationFrameCount)

	var wg sync.WaitGroup
	wg.Add(animationFrameCount)

	for i := range animationFrameCount {
		go func(i int) {
			defer wg.Done()

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

			canvas := createCanvas(width, height)

			resizedImg := resizeImage(baseImg, int(float64(width)*scaleX), int(float64(height)*scaleY))
			pasteImage(canvas, resizedImg, offsetX, offsetY)

			petFrame := resizeImage(hands[i], width, height)
			pasteImage(canvas, petFrame, 0, 0)

			images[i] = canvas
		}(i)
	}

	wg.Wait()

	return images
}

func pasteImage(dest draw.Image, src image.Image, offsetX, offsetY int) {
	draw.Draw(dest, src.Bounds().Add(image.Pt(offsetX, offsetY)), src, image.Point{}, draw.Over)
}

// MakeAnimation generates and returns pet-pet animation frames.
func MakeAnimation(base image.Image, config Config) []image.Image {
	base = prepareAnimationBase(base, config)

	return createAnimationFrames(base, config, func(width, height int) draw.Image {
		return image.NewRGBA(image.Rect(0, 0, width, height))
	})
}
