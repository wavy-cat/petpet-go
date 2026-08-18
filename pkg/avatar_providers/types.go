package avatar_providers

import "context"

const HttpUserAgent = "PetPet-Go/1.0"

type UserAvatar interface {
	GetId(ctx context.Context) (string, error)
	GetImage(ctx context.Context) ([]byte, error) // must be a png
}

type Provider interface {
	GetUserAvatar(ctx context.Context, userId string) (UserAvatar, error)
}
