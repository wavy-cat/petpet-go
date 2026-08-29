package animation

import (
	"context"
	"image"
)

type Service interface {
	GetMimeContentType() string
	GetOrGenerate(ctx context.Context, userID string, delay int) ([]byte, error)
	GenerateFromImage(ctx context.Context, img image.Image, delay int) ([]byte, error)
}
