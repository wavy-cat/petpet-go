package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type errReadCloser struct{}

func (errReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func (errReadCloser) Close() error {
	return nil
}

func TestNewBotAndNewBotWithClient(t *testing.T) {
	t.Parallel()

	bot := NewBot("token")
	if bot.token != "token" {
		t.Fatalf("expected token to be stored, got %q", bot.token)
	}
	if bot.client != http.DefaultClient {
		t.Fatalf("expected NewBot to use http.DefaultClient")
	}

	client := &http.Client{}
	bot = NewBotWithClient("token-2", client)
	if bot.token != "token-2" {
		t.Fatalf("expected token to be stored, got %q", bot.token)
	}
	if bot.client != client {
		t.Fatalf("expected provided client to be used")
	}

	bot = NewBotWithClient("token-3", nil)
	if bot.client != http.DefaultClient {
		t.Fatalf("expected nil client to fall back to http.DefaultClient")
	}
}

func TestNewUserByID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var gotReq *http.Request
		bot := NewBotWithClient("abc", &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			gotReq = req
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"id":"123","avatar":"hash"}`)),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})

		user, err := bot.NewUserById(context.Background(), "123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.ID != "123" {
			t.Fatalf("expected ID 123, got %q", user.ID)
		}
		if user.Avatar == nil || *user.Avatar != "hash" {
			t.Fatalf("expected avatar hash, got %#v", user.Avatar)
		}
		if gotReq == nil {
			t.Fatal("expected request to be sent")
		}
		if gotReq.Method != http.MethodGet {
			t.Fatalf("expected GET request, got %s", gotReq.Method)
		}
		if gotReq.URL.String() != baseURL+"/users/123" {
			t.Fatalf("unexpected URL: %s", gotReq.URL.String())
		}
		if gotReq.Header.Get("Authorization") != "Bot abc" {
			t.Fatalf("unexpected Authorization header: %q", gotReq.Header.Get("Authorization"))
		}
		if gotReq.Header.Get("User-Agent") != userAgent {
			t.Fatalf("unexpected User-Agent header: %q", gotReq.Header.Get("User-Agent"))
		}
		if gotReq.Header.Get("Accept") != "application/json" {
			t.Fatalf("unexpected Accept header: %q", gotReq.Header.Get("Accept"))
		}
	})

	t.Run("request error", func(t *testing.T) {
		t.Parallel()

		bot := NewBotWithClient("abc", http.DefaultClient)
		_, err := bot.NewUserById(context.Background(), "bad\n")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("client error", func(t *testing.T) {
		t.Parallel()

		want := errors.New("boom")
		bot := NewBotWithClient("abc", &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, want
		})})
		_, err := bot.NewUserById(context.Background(), "123")
		if !errors.Is(err, want) {
			t.Fatalf("expected client error %v, got %v", want, err)
		}
	})

	t.Run("non-200 response", func(t *testing.T) {
		t.Parallel()

		bot := NewBotWithClient("abc", &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Body:       io.NopCloser(strings.NewReader("not found")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})
		_, err := bot.NewUserById(context.Background(), "123")
		if err == nil || err.Error() != "error: not found" {
			t.Fatalf("expected formatted error, got %v", err)
		}
	})

	t.Run("body read error", func(t *testing.T) {
		t.Parallel()

		bot := NewBotWithClient("abc", &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       errReadCloser{},
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})
		_, err := bot.NewUserById(context.Background(), "123")
		if err == nil || err.Error() != "read error" {
			t.Fatalf("expected read error, got %v", err)
		}
	})

	t.Run("unmarshal error", func(t *testing.T) {
		t.Parallel()

		bot := NewBotWithClient("abc", &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("not json")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})
		_, err := bot.NewUserById(context.Background(), "123")
		if err == nil {
			t.Fatal("expected unmarshal error")
		}
	})
}

func TestGetAvatar(t *testing.T) {
	oldTransport := http.DefaultClient.Transport
	t.Cleanup(func() { http.DefaultClient.Transport = oldTransport })

	t.Run("success via wrapper", func(t *testing.T) {
		http.DefaultClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				return nil, fmt.Errorf("unexpected method %s", req.Method)
			}
			if req.URL.String() != baseCDNURL+"/avatars/123/hash.png?size=128" {
				return nil, fmt.Errorf("unexpected URL %s", req.URL.String())
			}
			if req.Header.Get("User-Agent") != userAgent {
				return nil, fmt.Errorf("unexpected User-Agent %q", req.Header.Get("User-Agent"))
			}
			if req.Header.Get("Accept") != "image/png" {
				return nil, fmt.Errorf("unexpected Accept %q", req.Header.Get("Accept"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("avatar-bytes")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})

		avatar := "hash"
		user := User{ID: "123", Avatar: &avatar}
		got, err := user.GetAvatar(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(got) != "avatar-bytes" {
			t.Fatalf("unexpected body %q", string(got))
		}
	})

	t.Run("nil client fallback", func(t *testing.T) {
		http.DefaultClient.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("avatar-bytes")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})

		avatar := "hash"
		got, err := User{ID: "123", Avatar: &avatar}.GetAvatarWithClient(context.Background(), nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if string(got) != "avatar-bytes" {
			t.Fatalf("unexpected body %q", string(got))
		}
	})

	t.Run("nil avatar", func(t *testing.T) {
		t.Parallel()

		_, err := User{ID: "123"}.GetAvatarWithClient(context.Background(), http.DefaultClient)
		if err == nil || err.Error() != "avatar not found" {
			t.Fatalf("expected avatar not found, got %v", err)
		}
	})

	t.Run("request error", func(t *testing.T) {
		t.Parallel()

		avatar := "bad\n"
		_, err := User{ID: "123", Avatar: &avatar}.GetAvatarWithClient(context.Background(), http.DefaultClient)
		if err == nil {
			t.Fatal("expected request error")
		}
	})

	t.Run("client error", func(t *testing.T) {
		t.Parallel()

		avatar := "hash"
		want := errors.New("boom")
		_, err := User{ID: "123", Avatar: &avatar}.GetAvatarWithClient(context.Background(), &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, want
		})})
		if !errors.Is(err, want) {
			t.Fatalf("expected client error %v, got %v", want, err)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		t.Parallel()

		avatar := "hash"
		_, err := User{ID: "123", Avatar: &avatar}.GetAvatarWithClient(context.Background(), &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     "500 Internal Server Error",
				Body:       io.NopCloser(strings.NewReader("nope")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})
		if err == nil || err.Error() != "invalid response status:500 Internal Server Error" {
			t.Fatalf("expected invalid status error, got %v", err)
		}
	})

	t.Run("body read error", func(t *testing.T) {
		t.Parallel()

		avatar := "hash"
		_, err := User{ID: "123", Avatar: &avatar}.GetAvatarWithClient(context.Background(), &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       errReadCloser{},
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})
		if err == nil || err.Error() != "read error" {
			t.Fatalf("expected read error, got %v", err)
		}
	})

	t.Run("not modified", func(t *testing.T) {
		t.Parallel()

		avatar := "hash"
		got, err := User{ID: "123", Avatar: &avatar}.GetAvatarWithClient(context.Background(), &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotModified,
				Status:     "304 Not Modified",
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		})})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected empty body, got %q", string(got))
		}
	})
}
