package test

import (
	"bytes"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"io"
	"os"
	"testing"
)

func Test(t *testing.T) {
	rawSource, err := os.Open("tasica.png")
	if err != nil {
		t.Fatal(err)
	}
	defer rawSource.Close()

	source, err := io.ReadAll(rawSource)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Generate Gif", func(t *testing.T) {
		t.Parallel()

		output, err := os.Create("output.gif")
		if err != nil {
			t.Fatal(err)
		}
		defer output.Close()

		err = petpet.MakeGif(bytes.NewReader(source), output, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
		if err != nil {
			t.Fatal("Error:", err)
		}
	})

	t.Run("Generate apng", func(t *testing.T) {
		t.Parallel()

		output, err := os.Create("output.apng")
		if err != nil {
			t.Fatal(err)
		}
		defer output.Close()

		err = petpet.MakeAPNG(bytes.NewReader(source), output, petpet.DefaultConfig)
		if err != nil {
			t.Fatal("Error:", err)
		}
	})

	t.Run("Generate faster apng", func(t *testing.T) {
		t.Parallel()

		output, err := os.Create("output_fast.apng")
		if err != nil {
			t.Fatal(err)
		}
		defer output.Close()

		config := petpet.DefaultConfig
		config.Delay = 2
		err = petpet.MakeAPNG(bytes.NewReader(source), output, config)
		if err != nil {
			t.Fatal("Error:", err)
		}
	})
}
