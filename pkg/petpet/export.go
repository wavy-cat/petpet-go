package petpet

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"

	"github.com/HugoSmits86/nativewebp"
)

const (
	maxPaletteColors           = 256
	webpMillisecondsPerGIFTick = 10
)

// ExportGIF writes animation frames in GIF format.
func ExportGIF(w io.Writer, images []image.Image, delay int, disposal byte) error {
	if err := validateImages(images); err != nil {
		return err
	}

	palette, err := createPalette(true, colorCountedImage{
		Image:      images[0],
		ColorCount: maxPaletteColors - 1,
	})
	if err != nil {
		return err
	}

	return gif.EncodeAll(w, &gif.GIF{
		Image:    palettizeImages(images, palette),
		Delay:    repeatedInts(len(images), delay),
		Disposal: repeatedBytes(len(images), disposal),
	})
}

// ExportWebp writes animation frames in WebP format.
func ExportWebp(w io.Writer, images []image.Image, delay int, disposal byte) error {
	if err := validateImages(images); err != nil {
		return err
	}

	return nativewebp.EncodeAll(w, &nativewebp.Animation{
		Images:          images,
		Durations:       webpDurations(len(images), delay),
		Disposals:       webpDisposals(len(images), disposal),
		LoopCount:       0,
		BackgroundColor: 0,
	}, nil)
}

var errNoImages = errors.New("must provide at least one image")

func validateImages(images []image.Image) error {
	if len(images) == 0 {
		return errNoImages
	}

	return nil
}

func palettizeImages(images []image.Image, palette color.Palette) []*image.Paletted {
	palettedImages := make([]*image.Paletted, len(images))
	for i, src := range images {
		dest := image.NewPaletted(src.Bounds(), palette)
		draw.Draw(dest, dest.Bounds(), src, src.Bounds().Min, draw.Src)
		palettedImages[i] = dest
	}

	return palettedImages
}

func repeatedInts(length, value int) []int {
	values := make([]int, length)
	for i := range values {
		values[i] = value
	}

	return values
}

func repeatedBytes(length int, value byte) []byte {
	values := make([]byte, length)
	for i := range values {
		values[i] = value
	}

	return values
}

func webpDurations(frameCount, delay int) []uint {
	durations := make([]uint, frameCount)
	for i := range durations {
		durations[i] = uint(delay * webpMillisecondsPerGIFTick)
	}

	return durations
}

func webpDisposals(frameCount int, disposal byte) []uint {
	disposals := make([]uint, frameCount)
	if disposal == 0 {
		return disposals
	}

	for i := range disposals {
		disposals[i] = 1
	}

	return disposals
}
func createPalette(addTransparent bool, images ...colorCountedImage) (color.Palette, error) {
	palette := make([]color.Color, 0, maxPaletteColors)
	if addTransparent {
		palette = append(palette, color.RGBA{})
	}

	for _, val := range images {
		imgPalette, err := quantizeImage(val.Image, val.ColorCount)
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
