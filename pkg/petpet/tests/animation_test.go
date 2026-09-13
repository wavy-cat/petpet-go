package tests_test

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/petpet"
)

func TestAnimation(t *testing.T) {
	t.Parallel()

	images := []struct {
		img        image.Image
		webpLen    int
		webpLenTol float64
		gifLen     int
	}{
		{
			img:        getImage("wavycat.png"),
			webpLen:    35202,
			webpLenTol: 0.10,
			gifLen:     78768,
		},
		{
			img:        getImage("tasica.png"),
			webpLen:    50648,
			webpLenTol: 0.10,
			gifLen:     66969,
		},
	}

	t.Run("Generate WebP animation", func(t *testing.T) {
		t.Parallel()

		for _, img := range images {
			output := bytes.Buffer{}

			imgs := petpet.MakeAnimation(img.img, petpet.DefaultImageSize, petpet.DefaultImageSize)
			err := petpet.ExportWebp(&output, imgs, petpet.DefaultDelay, petpet.DefaultDisposal)
			if err != nil {
				t.Fatal("ExportWebp returned error:", err)
			}

			if output.Len() < int(float64(img.webpLen)*(1-img.webpLenTol)) ||
				output.Len() > int(float64(img.webpLen)*(1+img.webpLenTol)) {
				t.Fatalf("unexpected output length: got %d, want ~%d", output.Len(), img.webpLen)
			}
		}
	})

	t.Run("Generate GIF animation", func(t *testing.T) {
		t.Parallel()

		for _, img := range images {
			output := bytes.Buffer{}

			imgs := petpet.MakeAnimation(img.img, petpet.DefaultImageSize, petpet.DefaultImageSize)
			err := petpet.ExportGIF(&output, imgs, petpet.DefaultDelay, petpet.DefaultDisposal)
			if err != nil {
				t.Fatal("ExportGIF returned error:", err)
			}

			if output.Len() != img.gifLen {
				t.Fatalf("unexpected output length: got %d, want %d", output.Len(), img.gifLen)
			}
		}
	})
}

func getImage(filename string) image.Image {
	// Read the contents of the file
	rawSource, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = rawSource.Close()
	}()

	// Convert content from Reader type to []bytes
	source, err := io.ReadAll(rawSource)
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(bytes.NewReader(source))
	if err != nil {
		panic(err)
	}

	return img
}
