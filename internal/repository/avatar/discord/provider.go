package discord

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wavy-cat/petpet-go/internal/repository/avatar"
	"github.com/wavy-cat/petpet-go/pkg/discord"
)

type Provider struct {
	bot    *discord.Bot
	client *http.Client
}

func NewDiscordAvatarProvider(botToken string, client *http.Client) avatar.Provider {
	if client == nil {
		client = http.DefaultClient
	}

	return &Provider{
		bot:    discord.NewBotWithClient(botToken, client),
		client: client,
	}
}

func (p *Provider) GetUserAvatar(ctx context.Context, userId string) (avatar.UserAvatar, error) {
	user, err := p.bot.NewUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}

	return &userAvatar{user: user, client: p.client}, nil
}
