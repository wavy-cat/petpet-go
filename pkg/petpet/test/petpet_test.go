package test

import (
	"fmt"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"io"
	"os"
	"testing"
)

func Test(t *testing.T) {
	source, err := os.Open("tasica.png")
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()

	output, err := os.Create("output.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()

	t.Run("Add task", func(t *testing.T) {
		result, err := petpet.MakeGif(source, petpet.DefaultConfig)
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
