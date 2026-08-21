package cachefactory

import (
	"errors"
	"strings"
	"testing"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/pkg/cache"
)

func TestNewDisabledCache(t *testing.T) {
	t.Parallel()

	instance, err := New(config.Cache{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if _, err := instance.Pull("key"); !errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() error = %v, want %v", err, cache.ErrNotExists)
	}
	if err := instance.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestNewMemoryCache(t *testing.T) {
	t.Parallel()

	instance, err := New(config.Cache{
		Storage: "memory",
		Memory: config.CacheMemoryConfig{
			Capacity: 1,
		},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := instance.Push("key", []byte("value")); err != nil {
		t.Fatalf("Push() error = %v", err)
	}
	if _, err := instance.Pull("key"); err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
}

func TestNewUnknownStorage(t *testing.T) {
	t.Parallel()

	instance, err := New(config.Cache{Storage: "redis"})
	if err == nil {
		t.Fatal("New() error = nil, want unsupported storage error")
	}
	if instance != nil {
		t.Fatalf("New() cache = %T, want nil", instance)
	}
	if !strings.Contains(err.Error(), `unsupported cache storage "redis"`) {
		t.Fatalf("New() error = %q", err)
	}
}
