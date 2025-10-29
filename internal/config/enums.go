package config

type LoggerPreset string

const (
	ProdPreset LoggerPreset = "prod"
	DevPreset  LoggerPreset = "dev"
	GCPPreset  LoggerPreset = "gcp"
)
