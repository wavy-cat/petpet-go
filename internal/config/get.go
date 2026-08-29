package config

import (
	"errors"
	"io/fs"

	"github.com/ilyakaznacheev/cleanenv"
)

func GetYMLConfig(path string) (Config, error) {
	var cfg Config
	return cfg, cleanenv.ReadConfig(path, &cfg)
}

func GetEnvConfig() (Config, error) {
	var cfg Config
	return cfg, cleanenv.ReadEnv(&cfg)
}

func GetConfig(path string) (Config, error) {
	cfg, err := GetYMLConfig(path)
	if _, ok := errors.AsType[*fs.PathError](err); ok {
		cfg, err = GetEnvConfig()
	}
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
