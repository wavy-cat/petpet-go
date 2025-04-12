package tests

import (
	"bytes"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
)

func TestGIF(t *testing.T) {
	images := []struct {
		source []byte
		len    int
	}{
		{
			source: getSource("wavycat.png"),
			len:    77075,
		},
		{
			source: getSource("tasica.png"),
			len:    58173,
		},
	}

	t.Run("Generate Gif", func(t *testing.T) {
		t.Parallel()

		for _, img := range images {
			output := bytes.Buffer{}

			err := petpet.MakeGif(bytes.NewReader(img.source), &output, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
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

			err := petpet.MakeGif(bytes.NewReader(img.source), &output, config, quantizers.HierarhicalQuantizer{})
			if err != nil {
				t.Fatal("MakeGIF returned error:", err)
			}

			if output.Len() != img.len {
				t.Fatalf("unexpected output length: got %d, want %d", output.Len(), img.len)
			}
		}
	})
}
