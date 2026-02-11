package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wavy-cat/petpet-go/pkg/cache"
)

type FileSystemCache struct {
	cacheDir string
	ttl      time.Duration

	mu        sync.RWMutex
	closingWg sync.WaitGroup
	closed    bool

	cleanerStop chan struct{}
	cleanerDone chan struct{}
}

// NewFileSystemCache creates a new file system cache with the specified directory for storage.
func NewFileSystemCache(cacheDir string, ttl time.Duration) (*FileSystemCache, error) {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0o755)
		if err != nil {
			return nil, err
		}
	}

	fsc := &FileSystemCache{cacheDir: cacheDir, ttl: ttl}
	if ttl > 0 {
		fsc.cleanerStop = make(chan struct{})
		fsc.cleanerDone = make(chan struct{})
		go fsc.runCleaner()
	}

	return fsc, nil
}

func (fsc *FileSystemCache) Push(key string, value []byte) error {
	if err := fsc.beginOperation(); err != nil {
		return err
	}
	defer fsc.endOperation()

	filePath := filepath.Join(fsc.cacheDir, key)
	return os.WriteFile(filePath, value, 0o644)
}

func (fsc *FileSystemCache) Pull(key string) ([]byte, error) {
	if err := fsc.beginOperation(); err != nil {
		return nil, err
	}
	defer fsc.endOperation()

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
	fsc.mu.Lock()
	if fsc.closed {
		fsc.mu.Unlock()
		return fmt.Errorf("cache is already closed")
	}
	fsc.closed = true
	cleanerStop := fsc.cleanerStop
	cleanerDone := fsc.cleanerDone
	fsc.mu.Unlock()

	if cleanerStop != nil {
		close(cleanerStop)
		<-cleanerDone
	}

	fsc.closingWg.Wait()
	return nil
}

func (fsc *FileSystemCache) beginOperation() error {
	fsc.mu.RLock()
	defer fsc.mu.RUnlock()

	if fsc.closed {
		return fmt.Errorf("cache is closed")
	}

	fsc.closingWg.Add(1)
	return nil
}

func (fsc *FileSystemCache) endOperation() {
	fsc.closingWg.Done()
}

func (fsc *FileSystemCache) runCleaner() {
	defer close(fsc.cleanerDone)

	ticker := time.NewTicker(cleanupInterval(fsc.ttl))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fsc.cleanupExpired()
		case <-fsc.cleanerStop:
			return
		}
	}
}

func (fsc *FileSystemCache) cleanupExpired() {
	entries, err := os.ReadDir(fsc.cacheDir)
	if err != nil {
		return
	}

	expireBefore := time.Now().Add(-fsc.ttl)
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		info, err := entry.Info()
		if err != nil || info.ModTime().After(expireBefore) {
			continue
		}

		_ = os.Remove(filepath.Join(fsc.cacheDir, entry.Name()))
	}
}

func cleanupInterval(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return 0
	}
	if ttl > time.Minute {
		return time.Minute
	}
	return ttl
}
