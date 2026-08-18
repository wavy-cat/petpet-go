package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wavy-cat/petpet-go/pkg/avatar_providers"
)

type provider struct {
	token  string
	client *http.Client
}

func NewProvider(token string) avatar_providers.Provider {
	client := &http.Client{Transport: http.DefaultTransport.(*http.Transport).Clone()}
	return &provider{
		token:  token,
		client: client,
	}
}

func (p *provider) GetUserAvatar(ctx context.Context, userId string) (avatar_providers.UserAvatar, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, userId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", p.token))
	req.Header.Set("User-Agent", avatar_providers.HttpUserAgent)
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	var user struct {
		Avatar *string `json:"avatar"` // optional
	}
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	if user.Avatar == nil {
		return nil, errors.New("avatar not found")
	}

	return &userAvatar{id: *user.Avatar, userId: userId, client: p.client}, nil
}
