package fs

import (
	"errors"
	"testing"
	"time"

	"github.com/wavy-cat/petpet-go/pkg/cache"
)

func TestFileSystemCacheTTLRemovesExpiredFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fsc, err := NewFileSystemCache(dir, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileSystemCache() error = %v", err)
	}
	defer func() {
		_ = fsc.Close()
	}()

	if err := fsc.Push("expired-key", []byte("value")); err != nil {
		t.Fatalf("Push() error = %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		_, pullErr := fsc.Pull("expired-key")
		if errors.Is(pullErr, cache.ErrNotExists) {
			return
		}
		if pullErr != nil {
			t.Fatalf("Pull() unexpected error = %v", pullErr)
		}
		time.Sleep(25 * time.Millisecond)
	}

	t.Fatal("expected cached file to be deleted by TTL cleaner")
}

func TestFileSystemCacheClose(t *testing.T) {
	t.Parallel()

	fsc, err := NewFileSystemCache(t.TempDir(), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("NewFileSystemCache() error = %v", err)
	}

	if err := fsc.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	if err := fsc.Close(); err == nil {
		t.Fatal("expected error on second Close() call")
	}

	if err := fsc.Push("key", []byte("value")); err == nil {
		t.Fatal("expected Push() to fail on closed cache")
	}

	if _, err := fsc.Pull("key"); err == nil {
		t.Fatal("expected Pull() to fail on closed cache")
	}
}
