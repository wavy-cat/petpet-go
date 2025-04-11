package test

import (
	"bytes"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"io"
	"os"
	"testing"
)

func TestAPNG(t *testing.T) {
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

	t.Run("Generate APNG", func(t *testing.T) {
		t.Parallel()

		const bufferLen = 290783
		output := bytes.Buffer{}

		err = petpet.MakeAPNG(bytes.NewReader(source), &output, petpet.DefaultConfig)
		if err != nil {
			t.Fatal("MakeAPNG returned error:", err)
		}

		if output.Len() != bufferLen {
			t.Fatalf("unexpected output length: got %d, want %d", output.Len(), bufferLen)
		}
	})

	t.Run("Generate faster APNG", func(t *testing.T) {
		t.Parallel()

		const bufferLen = 290783
		output := bytes.Buffer{}

		config := petpet.DefaultConfig
		config.Delay = 2

		err = petpet.MakeAPNG(bytes.NewReader(source), &output, config)
		if err != nil {
			t.Fatal("MakeAPNG returned error:", err)
		}

		if output.Len() != bufferLen {
			t.Fatalf("unexpected output length: got %d, want %d", output.Len(), bufferLen)
		}
	})
}
