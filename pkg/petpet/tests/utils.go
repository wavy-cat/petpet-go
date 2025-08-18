package tests

import (
	"io"
	"os"
)

func getSource(filename string) []byte {
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

	return source
}
