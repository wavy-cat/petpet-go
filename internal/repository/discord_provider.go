package repository

import (
	"context"
	"errors"
	"github.com/wavy-cat/petpet-go/pkg/discord"
)

type DiscordAvatarProvider struct {
	bot *discord.Bot
}

func NewDiscordAvatarProvider(bot *discord.Bot) AvatarProvider {
	return &DiscordAvatarProvider{bot: bot}
}

func (d *DiscordAvatarProvider) GetAvatarId(ctx context.Context, userId string) (string, error) {
	// Create a user object
	user, err := d.bot.NewUserById(ctx, userId)
	if err != nil {
		return "", err
	}

	// Checking whether the user has an avatar
	if user.Avatar == nil {
		return "", errors.New("not exists")
	}

	return *user.Avatar, nil
}

func (d *DiscordAvatarProvider) GetAvatarImage(ctx context.Context, userId string) ([]byte, error) {
	// Create a user object
	user, err := d.bot.NewUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Checking whether the user has an avatar
	if user.Avatar == nil {
		return nil, errors.New("not exists")
	}

	// Return the user avatar
	return user.GetAvatar(ctx)
}
