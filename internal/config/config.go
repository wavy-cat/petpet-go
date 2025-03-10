package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"io/fs"
)

type Server struct {
	Host            string `yaml:"host" env:"HOST"`
	Port            uint16 `yaml:"port" env:"PORT" env-default:"3000"`
	ShutdownTimeout uint   `yaml:"shutdownTimeout" env:"SHUTDOWN_TIMEOUT" env-default:"5000"`
}

type Discord struct {
	BotToken string `yaml:"botToken" env:"BOT_TOKEN" env-required:"true"`
}

type Cache struct {
	Storage        string `yaml:"storage" env:"CACHE_STORAGE"`
	MemoryCapacity uint   `yaml:"memoryCapacity" env:"CACHE_MEMORY_CAPACITY" env-default:"100"`
	FSPath         string `yaml:"fsPath" env:"CACHE_FS_PATH" env-default:"./cache"`
}

type Config struct {
	Server
	Discord
	Cache
}

func GetConfig() (Config, error) {
	var cfg Config
	var pathError *fs.PathError

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if errors.As(err, &pathError) {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
