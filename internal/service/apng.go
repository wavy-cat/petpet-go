package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/wavy-cat/petpet-go/internal/repository"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/petpet"
	"go.uber.org/zap"
)

type APNGService interface {
	GetOrGenerateAPNG(ctx context.Context, userId, source string, delay int) ([]byte, error)
}

type apngService struct {
	config    petpet.Config
	cache     cache.BytesCache
	providers map[string]repository.AvatarProvider
}

func NewAPngService(cache cache.BytesCache, providers map[string]repository.AvatarProvider,
	config petpet.Config) APNGService {
	return &apngService{
		config:    config,
		cache:     cache,
		providers: providers,
	}
}

func (g apngService) GetOrGenerateAPNG(ctx context.Context, userId, source string, delay int) ([]byte, error) {
	// Getting the required provider
	provider, ok := g.providers[source]
	if !ok {
		return nil, errors.New("unknown avatar source")
	}

	// Getting the user's avatar id
	avatarId, err := provider.GetAvatarId(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get avatar error: %v", err)
	}

	// We check if the GIF is in the cache and if so, return it.
	cacheName := fmt.Sprintf("%s-%d", avatarId, delay)

	if g.cache != nil {
		cachedGif, err := g.cache.Pull(cacheName)
		if err == nil {
			return cachedGif, nil
		} else if err.Error() != "not exist" {
			logger, ok := ctx.Value("logger").(*zap.Logger)
			if ok {
				logger.Warn("Error when retrieving GIF from cache",
					zap.Error(err), zap.String("avatar_id", avatarId))
			}
		}
	}

	// Getting the user's avatar
	avatarImage, err := provider.GetAvatarImage(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Generating a GIF
	config := g.config
	config.Delay = delay
	avatarReader := bytes.NewReader(avatarImage)

	var buf bytes.Buffer
	defer buf.Reset()
	err = petpet.MakeAPNG(avatarReader, &buf, config)
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()

	// Add a GIF to the cache
	if g.cache != nil {
		go func() {
			_ = g.cache.Push(cacheName, data)
		}()
	}

	// Returning the result
	return data, nil
}
