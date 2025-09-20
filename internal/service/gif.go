package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/wavy-cat/petpet-go/internal/middleware"
	"github.com/wavy-cat/petpet-go/internal/repository/avatar"

	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"go.uber.org/zap"
)

type GIFService interface {
	GetOrGenerateGif(ctx context.Context, userId string, delay int) ([]byte, error)
}

type gifService struct {
	config    petpet.Config
	quantizer petpet.Quantizer
	cache     cache.BytesCache
	provider  avatar.Provider
}

func NewGIFService(cache cache.BytesCache, provider avatar.Provider,
	config petpet.Config, quantizer petpet.Quantizer) GIFService {
	return &gifService{
		config:    config,
		quantizer: quantizer,
		cache:     cache,
		provider:  provider,
	}
}

func (s gifService) GetOrGenerateGif(ctx context.Context, userId string, delay int) ([]byte, error) {
	// Getting the logger
	logger, ok := ctx.Value(middleware.LoggerKey).(*zap.Logger)
	if !ok {
		panic("missing logger in gif service")
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

	if s.cache != nil {
		cachedGif, err := s.cache.Pull(cacheName)
		if err == nil {
			return cachedGif, nil
		} else if !errors.Is(err, cache.ErrNotExists) {
			logger.Warn("Error when retrieving GIF from cache",
				zap.Error(err),
				zap.String("avatar_id", avatarId))
		}
	}

	// Getting the user's avatar
	avatarImage, err := s.getAvatarImage(ctx, userAvatar)
	if err != nil {
		return nil, err
	}

	// Generating a GIF
	config := s.config
	config.Delay = delay
	avatarReader := bytes.NewReader(avatarImage)

	var buf bytes.Buffer
	defer buf.Reset()
	err = petpet.MakeGif(avatarReader, &buf, config, s.quantizer)
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()

	// Add a GIF to the cache
	if s.cache != nil {
		go func() {
			_ = s.cache.Push(cacheName, data)
		}()
	}

	return data, nil
}

func (s gifService) getAvatarImage(ctx context.Context, userAvatar avatar.UserAvatar) ([]byte, error) {
	avatarId, err := userAvatar.GetId(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %v", err)
	}

	if s.cache != nil {
		cached, err := s.cache.Pull(fmt.Sprintf("avatar-%s", avatarId))
		if err == nil {
			return cached, nil
		}
	}

	avatarImage, err := userAvatar.GetImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving avatar id: %v", err)
	}

	if s.cache != nil {
		go func() {
			err = s.cache.Push(fmt.Sprintf("avatar-%s", avatarId), avatarImage)
			if err != nil {
				logger := ctx.Value(middleware.LoggerKey).(*zap.Logger)
				logger.Error("Error pulling avatar in cache",
					zap.Error(err),
					zap.String("avatar_id", avatarId))
			}
		}()
	}

	return avatarImage, nil
}
