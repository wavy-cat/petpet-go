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
func ExportGIF(w io.Writer, images []image.Image, config Config) error {
	if len(images) == 0 {
		return errors.New("must provide at least one image")
	}

	palette, err := createPalette(true, colorCountedImage{
		Image:      images[0],
		ColorCount: maxPaletteColors - 1,
	})
	if err != nil {
		return err
	}

	palettedImages := make([]*image.Paletted, len(images))
	for i, src := range images {
		dest := image.NewPaletted(src.Bounds(), palette)
		draw.Draw(dest, dest.Bounds(), src, src.Bounds().Min, draw.Src)
		palettedImages[i] = dest
	}

	delays := make([]int, len(images))
	disposals := make([]byte, len(images))
	for i := range images {
		delays[i] = config.Delay
		disposals[i] = config.Disposal
	}

	return gif.EncodeAll(w, &gif.GIF{
		Image:    palettedImages,
		Delay:    delays,
		Disposal: disposals,
	})
}

// ExportWebp writes animation frames in WebP format.
func ExportWebp(w io.Writer, images []image.Image, config Config) error {
	if len(images) == 0 {
		return errors.New("must provide at least one image")
	}

	durations := make([]uint, len(images))
	disposals := make([]uint, len(images))
	for i := range images {
		durations[i] = uint(config.Delay * webpMillisecondsPerGIFTick)
		if config.Disposal > 0 {
			disposals[i] = 1
		}
	}

	return nativewebp.EncodeAll(w, &nativewebp.Animation{
		Images:          images,
		Durations:       durations,
		Disposals:       disposals,
		LoopCount:       0,
		BackgroundColor: 0,
	}, nil)
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
