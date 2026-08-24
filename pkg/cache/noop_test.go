package cache_test

import (
	"errors"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/cache"
)

func TestNoop(t *testing.T) {
	t.Parallel()

	noop := cache.NewNoop()

	if err := noop.Push("key", []byte("value")); err != nil {
		t.Fatalf("Push() error = %v", err)
	}

	value, err := noop.Pull("key")
	if !errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() error = %v, want %v", err, cache.ErrNotExists)
	}
	if value != nil {
		t.Fatalf("Pull() value = %q, want nil", value)
	}

	if err := noop.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
