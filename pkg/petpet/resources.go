package petpet

import (
	"embed"
	"image"
	"image/gif"
	"io/fs"
)

//go:embed img/pet*.gif
var resources embed.FS

var hands = loadHands()

func loadHands() [animationFrameCount]image.Image {
	files, err := fs.Glob(resources, "img/*.gif")
	if err != nil {
		panic(err)
	}

	var loadedHands [animationFrameCount]image.Image
	for i, file := range files {
		reader, err := resources.Open(file)
		if err != nil {
			panic(err)
		}

		img, decodeErr := gif.Decode(reader)
		closeErr := reader.Close()
		if decodeErr != nil {
			panic(decodeErr)
		}
		if closeErr != nil {
			panic(closeErr)
		}

		loadedHands[i] = img
	}

	return loadedHands
}
