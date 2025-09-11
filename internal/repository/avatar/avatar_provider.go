package avatar

import "context"

type Provider interface {
	GetAvatarId(ctx context.Context, userId string) (string, error)       // The method returns the avatar ID
	GetAvatarImage(ctx context.Context, userId string) ([]byte, error)    // The method returns the avatar image
	GetAvatar(ctx context.Context, userId string) ([]byte, string, error) // The method returns the avatar image and its ID
}
