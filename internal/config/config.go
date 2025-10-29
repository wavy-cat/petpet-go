package config

import (
	"errors"
	"io/fs"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerHeartbeat struct {
	Enable bool   `yaml:"enable" env:"HEARTBEAT_ENABLE" env-default:"false"`
	Path   string `yaml:"path" env:"HEARTBEAT_PATH" env-default:"/ping"`
}

type ServerThrottle struct {
	Enable         bool `yaml:"enable" env:"THROTTLE_ENABLE" env-default:"false"`
	Limit          int  `yaml:"limit" env:"THROTTLE_LIMIT" env-default:"2"`
	Backlog        int  `yaml:"backlog" env:"THROTTLE_BACKLOG" env-default:"3"`
	BacklogTimeout uint `yaml:"backlogTimeout" env:"THROTTLE_BACKLOG_TIMEOUT" env-default:"5"` // in secs
}

type Server struct {
	Host            string          `yaml:"host" env:"HOST"`
	Port            uint16          `yaml:"port" env:"PORT" env-default:"3000"`
	ShutdownTimeout uint            `yaml:"shutdownTimeout" env:"SHUTDOWN_TIMEOUT" env-default:"5000"`
	Heartbeat       ServerHeartbeat `yaml:"heartbeat"`
	Throttle        ServerThrottle  `yaml:"throttle"`
}

type Discord struct {
	BotToken string `yaml:"botToken" env:"BOT_TOKEN" env-required:"true"`
}

type CacheMemoryConfig struct {
	Capacity uint `yaml:"capacity" env:"CACHE_MEMORY_CAPACITY" env-default:"100"`
}

type CacheFSConfig struct {
	Path string `yaml:"path" env:"CACHE_FS_PATH" env-default:"./cache"`
}

type CacheS3Config struct {
	Bucket    string `yaml:"bucket" env:"CACHE_S3_BUCKET"`
	Endpoint  string `yaml:"endpoint" env:"CACHE_S3_ENDPOINT"`
	Region    string `yaml:"region" env:"CACHE_S3_REGION" env-default:"us-east-1"`
	AccessKey string `yaml:"accessKey" env:"CACHE_S3_ACCESS_KEY"`
	SecretKey string `yaml:"secretKey" env:"CACHE_S3_SECRET_KEY"`
}

type Cache struct {
	Storage string            `yaml:"storage" env:"CACHE_STORAGE"`
	Memory  CacheMemoryConfig `yaml:"memory"`
	FS      CacheFSConfig     `yaml:"fs"`
	S3      CacheS3Config     `yaml:"s3"`
}

type Proxy struct {
	URL string `yaml:"url" env:"PROXY_URL"`
}

type Logger struct {
	Preset LoggerPreset `yaml:"preset" env:"LOGGER_PRESET"`
}

type Config struct {
	Server
	Discord
	Cache
	Proxy
	Logger
}

func GetYMLConfig(path string) (Config, error) {
	var cfg Config
	return cfg, cleanenv.ReadConfig(path, &cfg)
}

func GetEnvConfig() (Config, error) {
	var cfg Config
	return cfg, cleanenv.ReadEnv(&cfg)
}

func GetConfig(path string) (Config, error) {
	var pathError *fs.PathError

	cfg, err := GetYMLConfig(path)
	if errors.As(err, &pathError) {
		cfg, err = GetEnvConfig()
	}
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
