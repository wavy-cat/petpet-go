package memory_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/cache/memory"
)

func TestNewLRUCache(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(2)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}
	if lru == nil {
		t.Fatal("NewLRUCache() returned nil cache")
	}
}

func TestNewLRUCacheInvalidCapacity(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(0)
	if err == nil {
		t.Fatal("NewLRUCache() expected error for zero capacity")
	}
	if lru != nil {
		t.Fatalf("NewLRUCache() cache = %v, want nil", lru)
	}
}

func TestLRUCachePushAndPull(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(2)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	want := []byte("value")
	if err := lru.Push("key", want); err != nil {
		t.Fatalf("Push() error = %v", err)
	}

	got, err := lru.Pull("key")
	if err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("Pull() = %q, want %q", got, want)
	}
}

func TestLRUCachePullMissing(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(1)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	got, err := lru.Pull("missing")
	if !errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() error = %v, want %v", err, cache.ErrNotExists)
	}
	if got != nil {
		t.Fatalf("Pull() = %q, want nil", got)
	}
}

func TestLRUCachePushUpdatesExistingKey(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(2)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	if err := lru.Push("key", []byte("old")); err != nil {
		t.Fatalf("Push() old value error = %v", err)
	}
	if err := lru.Push("key", []byte("new")); err != nil {
		t.Fatalf("Push() new value error = %v", err)
	}

	got, err := lru.Pull("key")
	if err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
	if !bytes.Equal(got, []byte("new")) {
		t.Fatalf("Pull() = %q, want %q", got, []byte("new"))
	}
}

func TestLRUCacheEvictsLeastRecentlyUsed(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(2)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	if err := lru.Push("first", []byte("1")); err != nil {
		t.Fatalf("Push() first error = %v", err)
	}
	if err := lru.Push("second", []byte("2")); err != nil {
		t.Fatalf("Push() second error = %v", err)
	}
	if _, err := lru.Pull("first"); err != nil {
		t.Fatalf("Pull() first error = %v", err)
	}
	if err := lru.Push("third", []byte("3")); err != nil {
		t.Fatalf("Push() third error = %v", err)
	}

	if _, err := lru.Pull("second"); !errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() second error = %v, want %v", err, cache.ErrNotExists)
	}
	for key, want := range map[string][]byte{
		"first": []byte("1"),
		"third": []byte("3"),
	} {
		got, err := lru.Pull(key)
		if err != nil {
			t.Fatalf("Pull(%q) error = %v", key, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("Pull(%q) = %q, want %q", key, got, want)
		}
	}
}

func TestLRUCacheCapacityOneEvictsPreviousKey(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(1)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	if err := lru.Push("old", []byte("old-value")); err != nil {
		t.Fatalf("Push() old error = %v", err)
	}
	if err := lru.Push("new", []byte("new-value")); err != nil {
		t.Fatalf("Push() new error = %v", err)
	}

	if _, err := lru.Pull("old"); !errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() old error = %v, want %v", err, cache.ErrNotExists)
	}
	got, err := lru.Pull("new")
	if err != nil {
		t.Fatalf("Pull() new error = %v", err)
	}
	if !bytes.Equal(got, []byte("new-value")) {
		t.Fatalf("Pull() new = %q, want %q", got, []byte("new-value"))
	}
}

func TestLRUCacheHandlesEmptyKeyAndNilValue(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(1)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	if err := lru.Push("", nil); err != nil {
		t.Fatalf("Push() error = %v", err)
	}
	got, err := lru.Pull("")
	if err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
	if got != nil {
		t.Fatalf("Pull() = %q, want nil", got)
	}
}

func TestLRUCacheClose(t *testing.T) {
	t.Parallel()

	lru, err := memory.NewLRUCache(1)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	if err := lru.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := lru.Push("key", []byte("value")); err != nil {
		t.Fatalf("Push() after Close() error = %v", err)
	}
	if _, err := lru.Pull("key"); err != nil {
		t.Fatalf("Pull() after Close() error = %v", err)
	}
}
