package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Bot struct {
	token string // Secret authorization token
}

func NewBot(token string) *Bot {
	return &Bot{token: token}
}

func (b Bot) NewUserById(ctx context.Context, id string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+b.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
