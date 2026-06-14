package discord

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wavy-cat/petpet-go/pkg/discord"
)

type userAvatar struct {
	user   *discord.User
	client *http.Client
}

func (a userAvatar) GetId(_ context.Context) (string, error) {
	if a.user.Avatar == nil {
		return "", fmt.Errorf("avatar not exists")
	}

	return *a.user.Avatar, nil
}

func (a userAvatar) GetImage(ctx context.Context) ([]byte, error) {
	if a.user.Avatar == nil {
		return nil, fmt.Errorf("avatar not exists")
	}

	return a.user.GetAvatarWithClient(ctx, a.client)
}
