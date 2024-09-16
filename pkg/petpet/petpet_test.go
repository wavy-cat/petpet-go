package petpet

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	source, err := os.Open("source.png")
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()

	result, err := MakeGif(source, DefaultConfig)
	if err != nil {
		fmt.Println(err)
		t.Fatal("Error:", err)
	}

	output, err := os.Create("output.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()

	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatal(err)
	}

	_, err = output.Write(data)
	if err != nil {
		t.Fatal(err)
	}
}
