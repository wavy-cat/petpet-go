package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wavy-cat/petpet-go/pkg/avatarproviders"
)

type provider struct {
	token  string
	client *http.Client
}

type Option func(*provider)

// WithHTTPClient configures the HTTP client used for Discord API and CDN requests.
func WithHTTPClient(client *http.Client) Option {
	return func(provider *provider) {
		if client != nil {
			provider.client = client
		}
	}
}

func NewProvider(token string, options ...Option) avatarproviders.Provider {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &provider{token: token, client: &http.Client{}}
	}

	provider := &provider{
		token:  token,
		client: &http.Client{Transport: transport.Clone()},
	}
	for _, option := range options {
		option(provider)
	}
	return provider
}

func (p *provider) GetUserAvatar(ctx context.Context, userID string) (avatarproviders.UserAvatar, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", p.token))
	req.Header.Set("User-Agent", avatarproviders.HTTPUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	var user struct {
		Avatar *string `json:"avatar"` // optional
	}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if user.Avatar == nil {
		return nil, errors.New("avatar not found")
	}

	return &userAvatar{id: *user.Avatar, userID: userID, client: p.client}, nil
}
