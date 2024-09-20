package test

import (
	"fmt"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"github.com/wavy-cat/petpet-go/pkg/petpet/quantizers"
	"io"
	"os"
	"testing"
)

func Test(t *testing.T) {
	source, err := os.Open("tasica.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(source)

	output, err := os.Create("output.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(output)

	t.Run("Add task", func(t *testing.T) {
		result, err := petpet.MakeGif(source, petpet.DefaultConfig, quantizers.HierarhicalQuantizer{})
		if err != nil {
			fmt.Println(err)
			t.Fatal("Error:", err)
		}

		data, err := io.ReadAll(result)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = output.Write(data)
		if err != nil {
			t.Error(err)
			return
		}
	})
}
