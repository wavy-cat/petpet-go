package config

import "time"

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

type ServerTLS struct {
	Enable   bool   `yaml:"enable" env:"TLS_ENABLE" env-default:"false"`
	CertFile string `yaml:"certFile" env:"TLS_CERT_FILE"`
	KeyFile  string `yaml:"keyFile" env:"TLS_KEY_FILE"`
}

type Server struct {
	Host            string          `yaml:"host" env:"HOST"`
	Port            uint16          `yaml:"port" env:"PORT" env-default:"3000"`
	ShutdownTimeout uint            `yaml:"shutdownTimeout" env:"SHUTDOWN_TIMEOUT" env-default:"5000"`
	WriteTimeout    time.Duration   `yaml:"writeTimeout" env:"WRITE_TIMEOUT" env-default:"15s"`
	ReadTimeout     time.Duration   `yaml:"readTimeout" env:"READ_TIMEOUT" env-default:"15s"`
	IdleTimeout     time.Duration   `yaml:"idleTimeout" env:"IDLE_TIMEOUT" env-default:"60s"`
	Heartbeat       ServerHeartbeat `yaml:"heartbeat"`
	Throttle        ServerThrottle  `yaml:"throttle"`
	TLS             ServerTLS       `yaml:"tls"`
	EnableHTTP2     bool            `yaml:"enableHttp2" env:"ENABLE_HTTP2" env-default:"true"`
}

type Discord struct {
	Enable   bool   `yaml:"enable" env:"DISCORD_ENABLE" env-default:"true"`
	BotToken string `yaml:"botToken" env:"BOT_TOKEN"`
}

type CacheMemoryConfig struct {
	Capacity uint `yaml:"capacity" env:"CACHE_MEMORY_CAPACITY" env-default:"100"`
}

type CacheFSConfig struct {
	Path string `yaml:"path" env:"CACHE_FS_PATH" env-default:"./cache"`
	TTL  uint   `yaml:"ttl" env:"CACHE_FS_TTL"`
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

type CustomUpload struct {
	Enable        bool   `yaml:"enable" env:"CUSTOM_UPLOAD_ENABLE" env-default:"true"`
	MaxUploadSize uint64 `yaml:"maxUploadSize" env:"CUSTOM_MAX_UPLOAD_SIZE" env-default:"5242880"`
	MaxPixelCount uint   `yaml:"maxPixelCount" env:"CUSTOM_MAX_PIXEL_COUNT" env-default:"1000000"`
}

type Config struct {
	Server
	Discord
	Cache
	Proxy
	Logger
	CustomUpload `yaml:"customUpload"`
}
