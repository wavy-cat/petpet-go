package handler

import (
	"github.com/wavy-cat/petpet-go/pkg/discord"
)

func GetAvatarUsingBot(bot *discord.Bot, userId string) ([]byte, error) {
	user, err := bot.NewUserById(userId)
	if err != nil {
		return nil, err
	}
	avatarImage, err := user.GetAvatar()
	if err != nil {
		return nil, err
	}
	return avatarImage, nil
}
