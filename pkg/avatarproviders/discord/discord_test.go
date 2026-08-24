package discord_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/wavy-cat/petpet-go/pkg/avatarproviders"
	"github.com/wavy-cat/petpet-go/pkg/avatarproviders/discord"
)

const testUserID = "613651509015740416"

var (
	errCDNRequestFailed       = errors.New("CDN request failed")
	errReadFailed             = errors.New("read failed")
	errRequestFailed          = errors.New("request failed")
	errUnexpectedTransportUse = errors.New("transport must not be called")
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type errorReadCloser struct{}

func (errorReadCloser) Read([]byte) (int, error) {
	return 0, errReadFailed
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

func providerWithTransport(transport http.RoundTripper) avatarproviders.Provider {
	return discord.NewProvider("test-token", discord.WithHTTPClient(&http.Client{Transport: transport}))
}

func avatarWithImageTransport(t *testing.T, imageTransport http.RoundTripper) avatarproviders.UserAvatar {
	t.Helper()

	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if strings.Contains(request.URL.Path, "/users/") {
			return response(http.StatusOK, io.NopCloser(strings.NewReader(`{"avatar":"avatar-hash"}`))), nil
		}
		return imageTransport.RoundTrip(request)
	})
	avatar, err := providerWithTransport(transport).GetUserAvatar(context.Background(), testUserID)
	if err != nil {
		t.Fatalf("GetUserAvatar() error = %v", err)
	}
	return avatar
}

func TestNewProvider(t *testing.T) {
	t.Parallel()

	if provider := discord.NewProvider("test-token"); provider == nil {
		t.Fatal("NewProvider() returned nil")
	}
}

func TestProviderGetUserAvatar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ctx       context.Context
		transport http.RoundTripper
		wantID    string
		wantErr   func(error) bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				assertUserRequest(t, request)
				return response(http.StatusOK, io.NopCloser(strings.NewReader(`{"avatar":"avatar-hash"}`))), nil
			}),
			wantID: "avatar-hash",
		},
		{
			name: "request creation error",
			ctx:  nil,
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errUnexpectedTransportUse
			}),
			wantErr: anyError,
		},
		{
			name: "transport error",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errRequestFailed
			}),
			wantErr: isError(errRequestFailed),
		},
		{
			name: "body read error",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, errorReadCloser{}), nil
			}),
			wantErr: isError(errReadFailed),
		},
		{
			name: "non-200 status",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusUnauthorized, io.NopCloser(strings.NewReader("unauthorized"))), nil
			}),
			wantErr: hasMessage("error: unauthorized"),
		},
		{
			name: "invalid JSON",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, io.NopCloser(strings.NewReader("{"))), nil
			}),
			wantErr: hasPrefix("error retrieving user:"),
		},
		{
			name: "avatar missing",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, io.NopCloser(strings.NewReader(`{"avatar":null}`))), nil
			}),
			wantErr: hasMessage("avatar not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			avatar, err := providerWithTransport(test.transport).GetUserAvatar(test.ctx, testUserID)
			assertAvatarResult(t, avatar, err, test.wantID, test.wantErr)
		})
	}
}

func TestUserAvatarGetImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ctx       context.Context
		transport http.RoundTripper
		wantImage string
		wantErr   func(error) bool
	}{
		{
			name: "request creation error",
			ctx:  nil,
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errUnexpectedTransportUse
			}),
			wantErr: anyError,
		},
		{
			name: "transport error",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errCDNRequestFailed
			}),
			wantErr: isError(errCDNRequestFailed),
		},
		{
			name: "invalid status",
			ctx:  context.Background(),
			transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				assertImageRequest(t, request)
				return response(http.StatusNotFound, io.NopCloser(strings.NewReader("not found"))), nil
			}),
			wantErr: hasMessage("invalid response status:404 test status"),
		},
		{
			name: "body read error",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, errorReadCloser{}), nil
			}),
			wantErr: isError(errReadFailed),
		},
		{
			name: "not modified",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusNotModified, io.NopCloser(strings.NewReader("png bytes"))), nil
			}),
			wantImage: "png bytes",
		},
		{
			name: "success",
			ctx:  context.Background(),
			transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return response(http.StatusOK, io.NopCloser(strings.NewReader("png bytes"))), nil
			}),
			wantImage: "png bytes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			image, err := avatarWithImageTransport(t, test.transport).GetImage(test.ctx)
			assertImageResult(t, image, err, test.wantImage, test.wantErr)
		})
	}
}

func TestUserAvatarGetID(t *testing.T) {
	t.Parallel()

	avatar := avatarWithImageTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errUnexpectedTransportUse
	}))
	id, err := avatar.GetID(context.Background())
	if err != nil {
		t.Fatalf("GetID() error = %v", err)
	}
	if id != "avatar-hash" {
		t.Errorf("GetID() = %q, want %q", id, "avatar-hash")
	}
}

func assertUserRequest(t *testing.T, request *http.Request) {
	t.Helper()

	if request.Method != http.MethodGet {
		t.Errorf("method = %q, want %q", request.Method, http.MethodGet)
	}
	if got, want := request.URL.String(), "https://discord.com/api/v10/users/"+testUserID; got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
	if got, want := request.Header.Get("Authorization"), "Bot test-token"; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	if got, want := request.Header.Get("User-Agent"), avatarproviders.HTTPUserAgent; got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
	if got, want := request.Header.Get("Accept"), "application/json"; got != want {
		t.Errorf("Accept = %q, want %q", got, want)
	}
}

func assertImageRequest(t *testing.T, request *http.Request) {
	t.Helper()

	if got, want := request.URL.String(), "https://cdn.discordapp.com/avatars/"+testUserID+"/avatar-hash.png?size=128"; got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
	if got, want := request.Header.Get("User-Agent"), avatarproviders.HTTPUserAgent; got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
	if got, want := request.Header.Get("Accept"), "image/png"; got != want {
		t.Errorf("Accept = %q, want %q", got, want)
	}
}

func assertAvatarResult(
	t *testing.T,
	avatar avatarproviders.UserAvatar,
	err error,
	wantID string,
	wantErr func(error) bool,
) {
	t.Helper()

	if wantErr != nil {
		if avatar != nil {
			t.Fatalf("avatar = %#v, want nil", avatar)
		}
		if !wantErr(err) {
			t.Fatalf("error = %v, did not meet expectation", err)
		}
		return
	}
	if err != nil {
		t.Fatalf("GetUserAvatar() error = %v", err)
	}
	if avatar == nil {
		t.Fatal("GetUserAvatar() avatar = nil")
	}
	id, err := avatar.GetID(context.Background())
	if err != nil {
		t.Fatalf("GetID() error = %v", err)
	}
	if id != wantID {
		t.Errorf("GetID() = %q, want %q", id, wantID)
	}
}

func assertImageResult(
	t *testing.T,
	image []byte,
	err error,
	wantImage string,
	wantErr func(error) bool,
) {
	t.Helper()

	if wantErr != nil {
		if image != nil {
			t.Fatalf("image = %v, want nil", image)
		}
		if !wantErr(err) {
			t.Fatalf("error = %v, did not meet expectation", err)
		}
		return
	}
	if err != nil {
		t.Fatalf("GetImage() error = %v", err)
	}
	if string(image) != wantImage {
		t.Errorf("GetImage() = %q, want %q", image, wantImage)
	}
}

func anyError(err error) bool {
	return err != nil
}

func isError(want error) func(error) bool {
	return func(err error) bool {
		return errors.Is(err, want)
	}
}

func hasMessage(want string) func(error) bool {
	return func(err error) bool {
		return err != nil && err.Error() == want
	}
}

func hasPrefix(want string) func(error) bool {
	return func(err error) bool {
		return err != nil && strings.HasPrefix(err.Error(), want)
	}
}
