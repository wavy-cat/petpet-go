package cachefactory

import (
	"fmt"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/pkg/cache"
	"github.com/wavy-cat/petpet-go/pkg/cache/fs"
	"github.com/wavy-cat/petpet-go/pkg/cache/memory"
	"github.com/wavy-cat/petpet-go/pkg/cache/s3"
)

func New(cfg config.Cache) (cache.BytesCache, error) {
	switch cfg.Storage {
	case "":
		return cache.NewNoop(), nil
	case "memory":
		instance, err := memory.NewLRUCache(cfg.Memory.Capacity)
		if err != nil {
			return nil, fmt.Errorf("create memory cache: %w", err)
		}
		return instance, nil
	case "fs":
		instance, err := fs.NewFileSystemCache(cfg.FS.Path, cfg.FS.TTL)
		if err != nil {
			return nil, fmt.Errorf("create filesystem cache: %w", err)
		}
		return instance, nil
	case "s3":
		instance, err := s3.NewS3Cache(
			cfg.S3.Bucket,
			cfg.S3.Endpoint,
			cfg.S3.Region,
			cfg.S3.AccessKey,
			cfg.S3.SecretKey,
		)
		if err != nil {
			return nil, fmt.Errorf("create S3 cache: %w", err)
		}
		return instance, nil
	default:
		return nil, fmt.Errorf("unsupported cache storage %q", cfg.Storage)
	}
}
