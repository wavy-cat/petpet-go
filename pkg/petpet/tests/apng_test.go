package tests

import (
	"bytes"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/petpet"
)

func TestAPNG(t *testing.T) {
	images := []struct {
		source []byte
		len    int
	}{
		{
			source: getSource("wavycat.png"),
			len:    295101,
		},
		{
			source: getSource("tasica.png"),
			len:    290783,
		},
	}

	t.Run("Generate APNG", func(t *testing.T) {
		t.Parallel()

		for _, img := range images {
			output := bytes.Buffer{}

			err := petpet.MakeAPNG(bytes.NewReader(img.source), &output, petpet.DefaultConfig)
			if err != nil {
				t.Fatal("MakeAPNG returned error:", err)
			}

			if output.Len() != img.len {
				t.Fatalf("unexpected output length: got %d, want %d", output.Len(), img.len)
			}
		}
	})

	t.Run("Generate faster APNG", func(t *testing.T) {
		t.Parallel()

		config := petpet.DefaultConfig
		config.Delay = 2

		for _, img := range images {
			output := bytes.Buffer{}

			err := petpet.MakeAPNG(bytes.NewReader(img.source), &output, petpet.DefaultConfig)
			if err != nil {
				t.Fatal("MakeAPNG returned error:", err)
			}

			if output.Len() != img.len {
				t.Fatalf("unexpected output length: got %d, want %d", output.Len(), img.len)
			}
		}
	})
}
