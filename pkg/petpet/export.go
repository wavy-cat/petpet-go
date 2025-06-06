package petpet

import (
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
