package cache

import (
	"errors"
	"testing"
)

func TestNoop(t *testing.T) {
	t.Parallel()

	cache := NewNoop()

	if err := cache.Push("key", []byte("value")); err != nil {
		t.Fatalf("Push() error = %v", err)
	}

	value, err := cache.Pull("key")
	if !errors.Is(err, ErrNotExists) {
		t.Fatalf("Pull() error = %v, want %v", err, ErrNotExists)
	}
	if value != nil {
		t.Fatalf("Pull() value = %q, want nil", value)
	}

	if err := cache.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
