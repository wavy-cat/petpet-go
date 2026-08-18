package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/avatar_providers"
)

const testUserID = "613651509015740416"

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type errorReadCloser struct{}

func (errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (errorReadCloser) Close() error {
	return nil
}

func response(statusCode int, body io.ReadCloser) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d test status", statusCode),
		Header:     make(http.Header),
		Body:       body,
	}
}

func nilContext() context.Context {
	return nil
}

func providerWithTransport(t *testing.T, transport http.RoundTripper) *provider {
	t.Helper()

	p, ok := NewProvider("test-token").(*provider)
	if !ok {
		t.Fatal("NewProvider returned an unexpected provider implementation")
	}
	p.client.Transport = transport
	return p
}

func TestNewProvider(t *testing.T) {
	t.Parallel()

	p, ok := NewProvider("test-token").(*provider)
	if !ok {
		t.Fatal("NewProvider returned an unexpected provider implementation")
	}
	if p.token != "test-token" {
		t.Errorf("token = %q, want %q", p.token, "test-token")
	}
	if p.client == nil {
		t.Fatal("client is nil")
	}
	if p.client.Transport == http.DefaultTransport {
		t.Error("provider client must use a cloned transport")
	}
}

func TestProviderGetUserAvatar(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		p := providerWithTransport(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Errorf("method = %q, want %q", req.Method, http.MethodGet)
			}
			if req.URL.String() != baseURL+"/users/"+testUserID {
				t.Errorf("URL = %q, want %q", req.URL.String(), baseURL+"/users/"+testUserID)
			}
			if got := req.Header.Get("Authorization"); got != "Bot test-token" {
				t.Errorf("Authorization = %q, want %q", got, "Bot test-token")
			}
			if got := req.Header.Get("User-Agent"); got != avatar_providers.HttpUserAgent {
				t.Errorf("User-Agent = %q, want %q", got, avatar_providers.HttpUserAgent)
			}
			if got := req.Header.Get("Accept"); got != "application/json" {
				t.Errorf("Accept = %q, want %q", got, "application/json")
			}
			return response(http.StatusOK, io.NopCloser(strings.NewReader(`{"avatar":"avatar-hash"}`))), nil
		}))

		avatar, err := p.GetUserAvatar(context.Background(), testUserID)
		if err != nil {
			t.Fatalf("GetUserAvatar() error = %v", err)
		}
		if avatar == nil {
			t.Fatal("GetUserAvatar() returned a nil avatar")
		}

		id, err := avatar.GetId(context.Background())
		if err != nil {
			t.Fatalf("GetId() error = %v", err)
		}
		if id != "avatar-hash" {
			t.Errorf("GetId() = %q, want %q", id, "avatar-hash")
		}
	})

	t.Run("request creation error", func(t *testing.T) {
		t.Parallel()

		p := providerWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("transport must not be called when request creation fails")
			return nil, nil
		}))

		avatar, err := p.GetUserAvatar(nilContext(), testUserID)
		if err == nil {
			t.Fatal("GetUserAvatar() error = nil, want an error")
		}
		if avatar != nil {
			t.Fatalf("GetUserAvatar() avatar = %#v, want nil", avatar)
		}
	})

	t.Run("transport error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("request failed")
		p := providerWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, wantErr
		}))

		avatar, err := p.GetUserAvatar(context.Background(), testUserID)
		if !errors.Is(err, wantErr) {
			t.Errorf("GetUserAvatar() error = %v, want %v", err, wantErr)
		}
		if avatar != nil {
			t.Fatalf("GetUserAvatar() avatar = %#v, want nil", avatar)
		}
	})

	t.Run("body read error", func(t *testing.T) {
		t.Parallel()

		p := providerWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
			return response(http.StatusOK, errorReadCloser{}), nil
		}))

		avatar, err := p.GetUserAvatar(context.Background(), testUserID)
		if err == nil || err.Error() != "read failed" {
			t.Errorf("GetUserAvatar() error = %v, want read failed", err)
		}
		if avatar != nil {
			t.Fatalf("GetUserAvatar() avatar = %#v, want nil", avatar)
		}
	})

	t.Run("non-200 status", func(t *testing.T) {
		t.Parallel()

		p := providerWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
			return response(http.StatusUnauthorized, io.NopCloser(strings.NewReader("unauthorized"))), nil
		}))

		avatar, err := p.GetUserAvatar(context.Background(), testUserID)
		if err == nil || err.Error() != "error: unauthorized" {
			t.Errorf("GetUserAvatar() error = %v, want error: unauthorized", err)
		}
		if avatar != nil {
			t.Fatalf("GetUserAvatar() avatar = %#v, want nil", avatar)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()

		p := providerWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
			return response(http.StatusOK, io.NopCloser(strings.NewReader("{"))), nil
		}))

		avatar, err := p.GetUserAvatar(context.Background(), testUserID)
		if err == nil || !strings.HasPrefix(err.Error(), "error retrieving user:") {
			t.Errorf("GetUserAvatar() error = %v, want JSON decoding error", err)
		}
		if avatar != nil {
			t.Fatalf("GetUserAvatar() avatar = %#v, want nil", avatar)
		}
	})

	t.Run("avatar missing", func(t *testing.T) {
		t.Parallel()

		p := providerWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
			return response(http.StatusOK, io.NopCloser(strings.NewReader(`{"avatar":null}`))), nil
		}))

		avatar, err := p.GetUserAvatar(context.Background(), testUserID)
		if err == nil || err.Error() != "avatar not found" {
			t.Errorf("GetUserAvatar() error = %v, want avatar not found", err)
		}
		if avatar != nil {
			t.Fatalf("GetUserAvatar() avatar = %#v, want nil", avatar)
		}
	})
}

func TestUserAvatarGetImage(t *testing.T) {
	t.Parallel()

	newAvatar := func(transport http.RoundTripper) userAvatar {
		return userAvatar{
			id:     "avatar-hash",
			userId: testUserID,
			client: &http.Client{Transport: transport},
		}
	}

	t.Run("id", func(t *testing.T) {
		t.Parallel()

		avatar := newAvatar(roundTripFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("transport must not be called by GetId")
			return nil, nil
		}))

		id, err := avatar.GetId(context.Background())
		if err != nil {
			t.Fatalf("GetId() error = %v", err)
		}
		if id != "avatar-hash" {
			t.Errorf("GetId() = %q, want %q", id, "avatar-hash")
		}
	})

	t.Run("request creation error", func(t *testing.T) {
		t.Parallel()

		avatar := newAvatar(roundTripFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("transport must not be called when request creation fails")
			return nil, nil
		}))

		image, err := avatar.GetImage(nilContext())
		if err == nil {
			t.Fatal("GetImage() error = nil, want an error")
		}
		if image != nil {
			t.Fatalf("GetImage() image = %v, want nil", image)
		}
	})

	t.Run("transport error", func(t *testing.T) {
		t.Parallel()

		wantErr := errors.New("CDN request failed")
		avatar := newAvatar(roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, wantErr
		}))

		image, err := avatar.GetImage(context.Background())
		if !errors.Is(err, wantErr) {
			t.Errorf("GetImage() error = %v, want %v", err, wantErr)
		}
		if image != nil {
			t.Fatalf("GetImage() image = %v, want nil", image)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		t.Parallel()

		avatar := newAvatar(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != baseCDNURL+"/avatars/"+testUserID+"/avatar-hash.png?size=128" {
				t.Errorf("URL = %q, want %q", req.URL.String(), baseCDNURL+"/avatars/"+testUserID+"/avatar-hash.png?size=128")
			}
			if got := req.Header.Get("User-Agent"); got != avatar_providers.HttpUserAgent {
				t.Errorf("User-Agent = %q, want %q", got, avatar_providers.HttpUserAgent)
			}
			if got := req.Header.Get("Accept"); got != "image/png" {
				t.Errorf("Accept = %q, want %q", got, "image/png")
			}
			return response(http.StatusNotFound, io.NopCloser(strings.NewReader("not found"))), nil
		}))

		image, err := avatar.GetImage(context.Background())
		if err == nil || err.Error() != "invalid response status:404 test status" {
			t.Errorf("GetImage() error = %v, want invalid response status error", err)
		}
		if image != nil {
			t.Fatalf("GetImage() image = %v, want nil", image)
		}
	})

	t.Run("body read error", func(t *testing.T) {
		t.Parallel()

		avatar := newAvatar(roundTripFunc(func(*http.Request) (*http.Response, error) {
			return response(http.StatusOK, errorReadCloser{}), nil
		}))

		image, err := avatar.GetImage(context.Background())
		if err == nil || err.Error() != "read failed" {
			t.Errorf("GetImage() error = %v, want read failed", err)
		}
		if image != nil {
			t.Fatalf("GetImage() image = %v, want nil", image)
		}
	})

	for _, statusCode := range []int{http.StatusNotModified, http.StatusOK} {
		statusCode := statusCode
		t.Run(fmt.Sprintf("success status %d", statusCode), func(t *testing.T) {
			t.Parallel()

			wantImage := []byte("png bytes")
			avatar := newAvatar(roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(statusCode, io.NopCloser(strings.NewReader(string(wantImage)))), nil
			}))

			image, err := avatar.GetImage(context.Background())
			if err != nil {
				t.Fatalf("GetImage() error = %v", err)
			}
			if string(image) != string(wantImage) {
				t.Errorf("GetImage() = %q, want %q", image, wantImage)
			}
		})
	}
}

var _ avatar_providers.Provider = (*provider)(nil)
var _ avatar_providers.UserAvatar = userAvatar{}
