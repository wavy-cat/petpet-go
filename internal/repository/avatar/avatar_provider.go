package avatar

import "context"

type UserAvatar interface {
	GetId(ctx context.Context) (string, error)
	GetImage(ctx context.Context) ([]byte, error)
}

type Provider interface {
	GetUserAvatar(ctx context.Context, userId string) (UserAvatar, error)
}
