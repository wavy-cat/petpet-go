package avatarproviders

import "context"

const HTTPUserAgent = "PetPet-Go/1.0"

type UserAvatar interface {
	GetID(ctx context.Context) (string, error)
	GetImage(ctx context.Context) ([]byte, error) // must be a png
}

type Provider interface {
	GetUserAvatar(ctx context.Context, userID string) (UserAvatar, error)
}
