package test

import (
	"bytes"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"io"
	"os"
	"testing"
)

func TestGIF(t *testing.T) {
	// Read the contents of the file
	rawSource, err := os.Open("tasica.png")
	if err != nil {
		t.Fatal(err)
	}

	// Convert content from Reader type to []bytes
	source, err := io.ReadAll(rawSource)
	if err != nil {
		t.Fatal(err)
	}

	// Closing the file
	err = rawSource.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Generate Gif", func(t *testing.T) {
		t.Parallel()

		const bufferLen = 58173
		output := bytes.Buffer{}

		err := petpet.MakeGif(bytes.NewReader(source), &output, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
		if err != nil {
			t.Fatal("generation error:", err)
		}

		if output.Len() != bufferLen {
			t.Fatalf("unexpected output length: got %d, want %d", output.Len(), bufferLen)
		}
	})

	t.Run("Generate faster Gif", func(t *testing.T) {
		t.Parallel()

		const bufferLen = 58173
		output := bytes.Buffer{}

		config := petpet.DefaultConfig
		config.Delay = 2

		err := petpet.MakeGif(bytes.NewReader(source), &output, config, quantizers.HierarhicalQuantizer{})
		if err != nil {
			t.Fatal("generation error:", err)
		}

		if output.Len() != bufferLen {
			t.Fatalf("unexpected output length: got %d, want %d", output.Len(), bufferLen)
		}
	})
}
