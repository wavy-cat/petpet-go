package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"

	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/pkg/avatarproviders"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"go.uber.org/zap"
)

type GIFService interface {
	GetOrGenerateGif(ctx context.Context, userID string, delay int) ([]byte, error)
	GenerateGifFromImage(ctx context.Context, img image.Image, delay int) ([]byte, error)
}

type gifService struct {
	cache    cache.BytesCache
	provider avatarproviders.Provider
}

func NewGIFService(cacheInstance cache.BytesCache, provider avatarproviders.Provider) GIFService {
	return &gifService{
		cache:    cacheInstance,
		provider: provider,
	}
}

func (s gifService) GetOrGenerateGif(ctx context.Context, userID string, delay int) ([]byte, error) {
	logger := middleware.LoggerFromContext(ctx)

	userAvatar, err := s.provider.GetUserAvatar(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar: %w", err)
	}

	avatarID, err := userAvatar.GetID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %w", err)
	}

	cacheName := fmt.Sprintf("%s-%d-gif", avatarID, delay)
	cachedGIF, err := s.cache.Pull(cacheName)
	if err == nil {
		return cachedGIF, nil
	}
	if !errors.Is(err, cache.ErrNotExists) {
		logger.Error("Error when pulling GIF from cache",
			zap.Error(err),
			zap.String("avatar_id", avatarID))
	}

	avatarImage, err := s.getAvatarImage(ctx, userAvatar)
	if err != nil {
		return nil, err
	}

	decodedImage, err := png.Decode(bytes.NewReader(avatarImage))
	if err != nil {
		return nil, fmt.Errorf("error decoding avatar image: %w", err)
	}

	data, err := s.GenerateGifFromImage(ctx, decodedImage, delay)
	if err != nil {
		return nil, err
	}

	go func(cacheName string, data []byte) {
		if err := s.cache.Push(cacheName, data); err != nil {
			logger.Error("Error when pushing GIF to cache",
				zap.Error(err),
				zap.String("avatar_id", avatarID))
		}
	}(cacheName, data)

	return data, nil
}

func (s gifService) GenerateGifFromImage(_ context.Context, img image.Image, delay int) ([]byte, error) {
	var buf bytes.Buffer
	defer buf.Reset()

	frames := petpet.MakeAnimation(img, petpet.DefaultImageSize, petpet.DefaultImageSize)
	if err := petpet.ExportGIF(&buf, frames, delay, petpet.DefaultDisposal); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s gifService) getAvatarImage(ctx context.Context, userAvatar avatarproviders.UserAvatar) ([]byte, error) {
	logger := middleware.LoggerFromContext(ctx)
	avatarID, err := userAvatar.GetID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %w", err)
	}

	cacheName := fmt.Sprintf("avatar-%s", avatarID)
	cached, err := s.cache.Pull(cacheName)
	if err == nil {
		return cached, nil
	}
	if !errors.Is(err, cache.ErrNotExists) {
		logger.Error("Error when pulling avatar from cache",
			zap.Error(err),
			zap.String("avatar_id", avatarID))
	}

	avatarImage, err := userAvatar.GetImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %w", err)
	}

	go func(cacheName string, avatarImage []byte) {
		if err := s.cache.Push(cacheName, avatarImage); err != nil {
			logger.Error("Error when pushing avatar to cache",
				zap.Error(err),
				zap.String("avatar_id", avatarID))
		}
	}(cacheName, avatarImage)

	return avatarImage, nil
}
