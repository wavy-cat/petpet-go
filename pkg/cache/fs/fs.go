package fs

import (
	"errors"
	"os"
	"path/filepath"
)

type FileSystemCache struct {
	cacheDir string
}

// NewFileSystemCache creates a new file system cache with the specified directory for storage.
func NewFileSystemCache(cacheDir string) (*FileSystemCache, error) {
	// Check if the cache directory exists, create it if not.
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return &FileSystemCache{cacheDir: cacheDir}, nil
}

// Push stores the data in the cache.
func (fsc *FileSystemCache) Push(key string, value []byte) error {
	// Generate the file path based on the key.
	filePath := filepath.Join(fsc.cacheDir, key)
	// Write the data to the file.
	return os.WriteFile(filePath, value, 0644)
}

// Pull retrieves the data from the cache.
func (fsc *FileSystemCache) Pull(key string) ([]byte, error) {
	// Generate the file path based on the key.
	filePath := filepath.Join(fsc.cacheDir, key)
	// Read the data from the file.
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("not exist")
		}
		return nil, err
	}
	return data, nil
}
