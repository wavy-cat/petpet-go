package loggerfactory

import (
	"fmt"

	"github.com/wavy-cat/petpet-go/internal/config"
	"github.com/wavy-cat/petpet-go/pkg/logger-presets/gcp"
	"go.uber.org/zap"
)

func New(cfg config.Logger) (*zap.Logger, error) {
	switch cfg.Preset {
	case "":
		return zap.NewNop(), nil
	case config.ProdPreset:
		return zap.NewProduction()
	case config.DevPreset:
		return zap.NewDevelopment()
	case config.GCPPreset:
		return gcp.NewGCPLogger()
	}

	return nil, fmt.Errorf("unsupported logger preset %q", cfg.Preset)
}

func NewWithService(cfg config.Logger, serviceName string) (*zap.Logger, error) {
	logger, err := New(cfg)
	if err != nil {
		return nil, err
	}

	return logger.With(zap.String("service", serviceName)), nil
}
