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
	baseImageScale        = 0.8
	horizontalScaleStep   = 0.02
	verticalScaleStep     = 0.05
)

func prepareAnimationBase(baseImg image.Image, width, height int) image.Image {
	if size := baseImg.Bounds().Size(); size.X != width || size.Y != height {
		return resizeImage(baseImg, width, height)
	}

	return baseImg
}

func createAnimationFrames(baseImg image.Image, width, height int) []image.Image {
	frames := make([]image.Image, animationFrameCount)

	var wg sync.WaitGroup
	wg.Add(len(frames))

	for frameIndex := range frames {
		go func(frameIndex int) {
			defer wg.Done()
			frames[frameIndex] = createAnimationFrame(baseImg, width, height, frameIndex)
		}(frameIndex)
	}

	wg.Wait()

	return frames
}

func createAnimationFrame(baseImg image.Image, width, height, frameIndex int) image.Image {
	scaleX, scaleY := animationFrameScale(frameIndex)
	offsetX := int(((1-scaleX)*0.5 + horizontalImageOffset) * float64(width))
	offsetY := int(((1 - scaleY) - verticalImageOffset) * float64(height))

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	resizedBase := resizeImage(
		baseImg,
		int(float64(width)*scaleX),
		int(float64(height)*scaleY),
	)
	pasteImage(canvas, resizedBase, offsetX, offsetY)
	pasteImage(canvas, resizeImage(hands[frameIndex], width, height), 0, 0)

	return canvas
}

func animationFrameScale(frameIndex int) (float64, float64) {
	progress := frameIndex
	if frameIndex >= animationFrameCount/2 {
		progress = animationFrameCount - frameIndex
	}

	return baseImageScale + float64(progress)*horizontalScaleStep,
		baseImageScale - float64(progress)*verticalScaleStep
}

func pasteImage(dest draw.Image, src image.Image, offsetX, offsetY int) {
	destinationMin := image.Pt(offsetX, offsetY)
	destinationBounds := image.Rectangle{
		Min: destinationMin,
		Max: destinationMin.Add(src.Bounds().Size()),
	}
	draw.Draw(dest, destinationBounds, src, src.Bounds().Min, draw.Over)
}

// MakeAnimation generates and returns pet-pet animation frames.
func MakeAnimation(base image.Image, width, height int) []image.Image {
	return createAnimationFrames(prepareAnimationBase(base, width, height), width, height)
}
