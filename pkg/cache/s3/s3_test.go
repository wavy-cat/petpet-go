package s3_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/cache"
	s3cache "github.com/wavy-cat/petpet-go/pkg/cache/s3"
)

func newTestCache(t *testing.T, handler http.HandlerFunc) *s3cache.Cache {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	s3Cache, err := s3cache.NewS3Cache("bucket", server.URL, "us-east-1", "test-access-key", "test-secret-key")
	if err != nil {
		t.Fatalf("NewS3Cache() error = %v", err)
	}

	return s3Cache
}

func TestNewS3Cache(t *testing.T) {
	t.Parallel()

	s3Cache, err := s3cache.NewS3Cache("bucket", "https://example.com", "us-west-2", "access", "secret")
	if err != nil {
		t.Fatalf("NewS3Cache() error = %v", err)
	}
	if s3Cache == nil {
		t.Fatal("NewS3Cache() returned nil cache")
	}
}

func TestNewS3CacheRequiresBucketName(t *testing.T) {
	t.Parallel()

	s3Cache, err := s3cache.NewS3Cache("", "https://example.com", "us-west-2", "access", "secret")
	if err == nil {
		t.Fatal("NewS3Cache() expected error for empty bucket name")
	}
	if s3Cache != nil {
		t.Fatalf("NewS3Cache() cache = %v, want nil", s3Cache)
	}
}

func TestCachePushAndPull(t *testing.T) {
	t.Parallel()

	s3Cache := newTestCache(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			if got, want := r.URL.Path, "/bucket/key"; got != want {
				t.Fatalf("PUT path = %q, want %q", got, want)
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}
			if !bytes.Equal(body, []byte("value")) {
				t.Fatalf("PUT body = %q, want %q", body, []byte("value"))
			}
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			if got, want := r.URL.Path, "/bucket/key"; got != want {
				t.Fatalf("GET path = %q, want %q", got, want)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("value")); err != nil {
				t.Errorf("Write() error = %v", err)
			}
		default:
			t.Fatalf("unexpected method %q", r.Method)
		}
	})

	if err := s3Cache.Push("key", []byte("value")); err != nil {
		t.Fatalf("Push() error = %v", err)
	}

	got, err := s3Cache.Pull("key")
	if err != nil {
		t.Fatalf("Pull() error = %v", err)
	}
	if !bytes.Equal(got, []byte("value")) {
		t.Fatalf("Pull() = %q, want %q", got, []byte("value"))
	}
}

func TestCachePushEmptyKeyAndNilValue(t *testing.T) {
	t.Parallel()

	s3Cache := newTestCache(t, func(_ http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request for invalid empty key: %s %s", r.Method, r.URL.Path)
	})

	if err := s3Cache.Push("", nil); err == nil {
		t.Fatal("Push() expected error for empty key")
	}
}

func TestCachePullMissingReturnsErrNotExists(t *testing.T) {
	t.Parallel()

	s3Cache := newTestCache(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want %q", r.Method, http.MethodGet)
		}
		w.WriteHeader(http.StatusNotFound)
	})

	got, err := s3Cache.Pull("missing")
	if !errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() error = %v, want %v", err, cache.ErrNotExists)
	}
	if got != nil {
		t.Fatalf("Pull() = %q, want nil", got)
	}
}

func TestCachePullPropagatesOtherErrors(t *testing.T) {
	t.Parallel()

	s3Cache := newTestCache(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want %q", r.Method, http.MethodGet)
		}
		w.WriteHeader(http.StatusInternalServerError)
	})

	got, err := s3Cache.Pull("broken")
	if err == nil {
		t.Fatal("Pull() expected error")
	}
	if errors.Is(err, cache.ErrNotExists) {
		t.Fatalf("Pull() error = %v, want non-ErrNotExists error", err)
	}
	if got != nil {
		t.Fatalf("Pull() = %q, want nil", got)
	}
}

func TestCacheClose(t *testing.T) {
	t.Parallel()

	s3Cache := newTestCache(t, func(_ http.ResponseWriter, _ *http.Request) {})

	if err := s3Cache.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := s3Cache.Close(); err == nil {
		t.Fatal("Close() expected error on second call")
	}
	if err := s3Cache.Push("key", []byte("value")); err == nil {
		t.Fatal("Push() expected error after Close()")
	}
	if _, err := s3Cache.Pull("key"); err == nil {
		t.Fatal("Pull() expected error after Close()")
	}
}
