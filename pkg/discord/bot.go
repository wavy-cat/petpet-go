package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Bot struct {
	token  string // Secret authorization token
	client *http.Client
}

func NewBot(token string) *Bot {
	return NewBotWithClient(token, http.DefaultClient)
}

func NewBotWithClient(token string, client *http.Client) *Bot {
	if client == nil {
		client = http.DefaultClient
	}

	return &Bot{token: token, client: client}
}

func (b Bot) NewUserById(ctx context.Context, id string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+b.token)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s", string(body))
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
