package memory

import (
	"bytes"
	"errors"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/cache"
)

func TestNewLRUCache(t *testing.T) {
	t.Parallel()

	lru, err := NewLRUCache(2)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}
	if lru == nil {
		t.Fatal("NewLRUCache() returned nil cache")
	}
	if lru.capacity != 2 {
		t.Fatalf("NewLRUCache() capacity = %d, want 2", lru.capacity)
	}
	if lru.cache == nil {
		t.Fatal("NewLRUCache() cache map is nil")
	}
	if lru.ll == nil {
		t.Fatal("NewLRUCache() list is nil")
	}
}

func TestNewLRUCacheInvalidCapacity(t *testing.T) {
	t.Parallel()

	lru, err := NewLRUCache(0)
	if err == nil {
		t.Fatal("NewLRUCache() expected error for zero capacity")
	}
	if lru != nil {
		t.Fatalf("NewLRUCache() cache = %v, want nil", lru)
	}
}

func TestLRUCachePushAndPull(t *testing.T) {
	t.Parallel()

	lru, err := NewLRUCache(2)
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

	lru, err := NewLRUCache(1)
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

	lru, err := NewLRUCache(2)
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
	if len(lru.cache) != 1 {
		t.Fatalf("cache size = %d, want 1", len(lru.cache))
	}
}

func TestLRUCacheEvictsLeastRecentlyUsed(t *testing.T) {
	t.Parallel()

	lru, err := NewLRUCache(2)
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

	lru, err := NewLRUCache(1)
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

	lru, err := NewLRUCache(1)
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

func TestLRUCacheRemoveOldestEmptyCache(t *testing.T) {
	t.Parallel()

	lru, err := NewLRUCache(1)
	if err != nil {
		t.Fatalf("NewLRUCache() error = %v", err)
	}

	lru.removeOldest()

	if len(lru.cache) != 0 {
		t.Fatalf("cache size = %d, want 0", len(lru.cache))
	}
	if lru.ll.Len() != 0 {
		t.Fatalf("list length = %d, want 0", lru.ll.Len())
	}
}

func TestLRUCacheClose(t *testing.T) {
	t.Parallel()

	lru, err := NewLRUCache(1)
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
