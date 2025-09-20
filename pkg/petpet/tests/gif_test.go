package tests

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
)

func TestGIF(t *testing.T) {
	images := []struct {
		img image.Image
		len int
	}{
		{
			img: getImage("wavycat.png"),
			len: 73599,
		},
		{
			img: getImage("tasica.png"),
			len: 55654,
		},
	}

	t.Run("Generate Gif", func(t *testing.T) {
		t.Parallel()

		for _, img := range images {
			output := bytes.Buffer{}

			err := petpet.MakeGif(img.img, &output, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
			if err != nil {
				t.Fatal("MakeGIF returned error:", err)
			}

			if output.Len() != img.len {
				t.Fatalf("unexpected output length: got %d, want %d", output.Len(), img.len)
			}
		}
	})

	t.Run("Generate faster Gif", func(t *testing.T) {
		t.Parallel()

		config := petpet.DefaultConfig
		config.Delay = 2

		for _, img := range images {
			output := bytes.Buffer{}

			err := petpet.MakeGif(img.img, &output, config, quantizers.HierarhicalQuantizer{})
			if err != nil {
				t.Fatal("MakeGIF returned error:", err)
			}

			if output.Len() != img.len {
				t.Fatalf("unexpected output length: got %d, want %d", output.Len(), img.len)
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
