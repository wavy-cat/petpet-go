package petpet

import (
	"embed"
	"fmt"
	"image"
	"image/gif"
	"io/fs"
)

//go:embed img/pet*.gif
var resources embed.FS

var hands [10]image.Image // Массив с GIF изображениями рук

func init() {
	files, err := fs.Glob(resources, "img/*.gif")
	if err != nil {
		panic(err)
	}

	for i, file := range files {
		reader, err := resources.Open(file)
		if err != nil {
			panic(err)
		}

		img, err := gif.Decode(reader)
		if err != nil {
			panic(err)
		}

		err = reader.Close()
		if err != nil {
			fmt.Println("Failed to close file reader")
		}

		hands[i] = img
	}
}
