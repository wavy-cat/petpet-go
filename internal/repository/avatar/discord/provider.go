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

func NewDiscordAvatarProvider(botToken string) avatar.Provider {
	return &Provider{bot: discord.NewBot(botToken)}
}

func (p *Provider) GetUserAvatar(ctx context.Context, userId string) (avatar.UserAvatar, error) {
	user, err := p.bot.NewUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	return &userAvatar{user: user}, nil
}
