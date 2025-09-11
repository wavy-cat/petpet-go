package discord

import (
	"context"
	"fmt"

	"github.com/wavy-cat/petpet-go/internal/repository/avatar"
	"github.com/wavy-cat/petpet-go/pkg/discord"
)

type Provider struct {
	bot *discord.Bot
}

func NewDiscordAvatarProvider(bot *discord.Bot) avatar.Provider {
	return &Provider{bot: bot}
}

func (d *Provider) GetAvatarId(ctx context.Context, userId string) (string, error) {
	user, err := d.bot.NewUserById(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("error retrieving user: %v", err)
	}

	if user.Avatar == nil {
		return "", fmt.Errorf("avatar not exists")
	}

	return *user.Avatar, nil
}

func (d *Provider) GetAvatarImage(ctx context.Context, userId string) ([]byte, error) {
	user, err := d.bot.NewUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	if user.Avatar == nil {
		return nil, fmt.Errorf("avatar not exists")
	}

	return user.GetAvatar(ctx)
}

func (d *Provider) GetAvatar(ctx context.Context, userId string) ([]byte, string, error) {
	user, err := d.bot.NewUserById(ctx, userId)
	if err != nil {
		return nil, "", fmt.Errorf("error retrieving user: %v", err)
	}

	if user.Avatar == nil {
		return nil, "", fmt.Errorf("avatar not exists")
	}

	avatarImg, err := user.GetAvatar(ctx)
	if user.Avatar == nil {
		return nil, "", fmt.Errorf("error while retrieving avatar: %v", err)
	}

	return avatarImg, *user.Avatar, nil
}
