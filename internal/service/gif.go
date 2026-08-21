package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"

	"github.com/wavy-cat/petpet-go/internal/middleware"

	"github.com/wavy-cat/petpet-go/pkg/avatar_providers"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"go.uber.org/zap"
)

type GIFService interface {
	GetOrGenerateGif(ctx context.Context, userId string, delay int) ([]byte, error)
	GenerateGifFromImage(ctx context.Context, img image.Image, delay int) ([]byte, error)
}

type gifService struct {
	config    petpet.Config
	quantizer petpet.Quantizer
	cache     cache.BytesCache
	provider  avatar_providers.Provider
}

func NewGIFService(cacheInstance cache.BytesCache, provider avatar_providers.Provider,
	config petpet.Config, quantizer petpet.Quantizer) GIFService {
	if cacheInstance == nil {
		cacheInstance = cache.NewNoop()
	}

	return &gifService{
		config:    config,
		quantizer: quantizer,
		cache:     cacheInstance,
		provider:  provider,
	}
}

func (s gifService) GetOrGenerateGif(ctx context.Context, userId string, delay int) ([]byte, error) {
	// Getting the logger-presets
	logger, ok := ctx.Value(middleware.LoggerKey).(*zap.Logger)
	if !ok {
		panic("missing logger-presets in gif service")
	}

	// Getting the user's avatar id
	userAvatar, err := s.provider.GetUserAvatar(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar: %v", err)
	}

	avatarId, err := userAvatar.GetId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %v", err)
	}

	// We check if the GIF is in the cache and if so, return it.
	cacheName := fmt.Sprintf("%s-%d-gif", avatarId, delay)

	cachedGif, err := s.cache.Pull(cacheName)
	if err == nil {
		return cachedGif, nil
	}
	if !errors.Is(err, cache.ErrNotExists) {
		logger.Error("Error when pulling GIF from cache",
			zap.Error(err),
			zap.String("avatar_id", avatarId))
	}

	// Getting the user's avatar
	avatarImage, err := s.getAvatarImage(ctx, userAvatar)
	if err != nil {
		return nil, err
	}

	decodedImage, err := png.Decode(bytes.NewReader(avatarImage))
	if err != nil {
		return nil, fmt.Errorf("error decoding avatar image: %v", err)
	}

	data, err := s.GenerateGifFromImage(ctx, decodedImage, delay)
	if err != nil {
		return nil, err
	}

	// Add a GIF to the cache
	go func(cacheName string, data []byte) {
		if err := s.cache.Push(cacheName, data); err != nil {
			logger.Error("Error when pushing GIF to cache",
				zap.Error(err),
				zap.String("avatar_id", avatarId))
		}
	}(cacheName, data)

	return data, nil
}

func (s gifService) GenerateGifFromImage(ctx context.Context, img image.Image, delay int) ([]byte, error) {
	_ = ctx
	config := s.config
	config.Delay = delay

	var buf bytes.Buffer
	defer buf.Reset()

	if err := petpet.MakeGif(img, &buf, config, s.quantizer); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s gifService) getAvatarImage(ctx context.Context, userAvatar avatar_providers.UserAvatar) ([]byte, error) {
	avatarId, err := userAvatar.GetId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %v", err)
	}

	cacheName := fmt.Sprintf("avatar-%s", avatarId)
	cached, err := s.cache.Pull(cacheName)
	if err == nil {
		return cached, nil
	}
	if !errors.Is(err, cache.ErrNotExists) {
		logger := ctx.Value(middleware.LoggerKey).(*zap.Logger)
		logger.Error("Error when pulling avatar from cache",
			zap.Error(err),
			zap.String("avatar_id", avatarId))
	}

	avatarImage, err := userAvatar.GetImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %v", err)
	}

	logger := ctx.Value(middleware.LoggerKey).(*zap.Logger)
	go func(cacheName string, avatarImage []byte) {
		if err := s.cache.Push(cacheName, avatarImage); err != nil {
			logger.Error("Error when pushing avatar to cache",
				zap.Error(err),
				zap.String("avatar_id", avatarId))
		}
	}(cacheName, avatarImage)

	return avatarImage, nil
}
