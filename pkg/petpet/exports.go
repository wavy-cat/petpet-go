package petpet

import (
	"github.com/kettek/apng"
	"image"
	"image/gif"
	"io"
)

func exportGIF(w io.Writer, images []*image.Paletted, delays []int, disposals []byte) error {
	g := &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposals,
	}

	return gif.EncodeAll(w, g)
}

func exportAPNG(w io.Writer, images []*image.RGBA, delay uint16) error {
	a := apng.APNG{
		Frames: make([]apng.Frame, len(images)),
	}

	for i, m := range images {
		a.Frames[i].Image = m
		a.Frames[i].DelayNumerator = delay
	}

	return apng.Encode(w, a)
}
