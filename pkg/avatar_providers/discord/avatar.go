package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wavy-cat/petpet-go/pkg/avatar_providers"
)

type userAvatar struct {
	id     string
	userId string
	client *http.Client
}

func (a userAvatar) GetId(_ context.Context) (string, error) {
	return a.id, nil
}

func (a userAvatar) GetImage(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/avatars/%s/%s.png?size=128", baseCDNURL, a.userId, a.id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", avatar_providers.HttpUserAgent)
	req.Header.Set("Accept", "image/png")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusNotModified && resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response status:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
