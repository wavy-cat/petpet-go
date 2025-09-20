package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wavy-cat/petpet-go/pkg/cache"
)

type FileSystemCache struct {
	cacheDir  string
	closingWg sync.WaitGroup
	closed    bool
}

// NewFileSystemCache creates a new file system cache with the specified directory for storage.
func NewFileSystemCache(cacheDir string) (*FileSystemCache, error) {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return &FileSystemCache{cacheDir: cacheDir}, nil
}

func (fsc *FileSystemCache) Push(key string, value []byte) error {
	if fsc.closed {
		return fmt.Errorf("cache is closed")
	}
	fsc.closingWg.Add(1)
	defer fsc.closingWg.Done()

	filePath := filepath.Join(fsc.cacheDir, key)
	return os.WriteFile(filePath, value, 0644)
}

func (fsc *FileSystemCache) Pull(key string) ([]byte, error) {
	if fsc.closed {
		return nil, fmt.Errorf("cache is closed")
	}
	fsc.closingWg.Add(1)
	defer fsc.closingWg.Done()

	filePath := filepath.Join(fsc.cacheDir, key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, cache.ErrNotExists
		}
		return nil, err
	}
	return data, nil
}

func (fsc *FileSystemCache) Close() error {
	if fsc.closed {
		return fmt.Errorf("cache is already closed")
	}
	fsc.closed = true
	fsc.closingWg.Wait()
	return nil
}
