package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	Avatar        *string `json:"avatar"`
}

func (u User) GetAvatar(ctx context.Context) ([]byte, error) {
	if u.Avatar == nil {
		return nil, errors.New("avatar not found")
	}

	url := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.ID, *u.Avatar)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotModified && resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response status:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
