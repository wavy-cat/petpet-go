package repository

import "context"

type AvatarProvider interface {
	GetAvatarId(ctx context.Context, userId string) (string, error)
	GetAvatarImage(ctx context.Context, userId string) ([]byte, error)
}
