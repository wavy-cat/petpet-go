package config

import "os"

var HTTPAddress = ":80"

func init() {
	if env := os.Getenv("ADDRESS"); env != "" {
		HTTPAddress = env
	} else if env := os.Getenv("PORT"); env != "" {
		HTTPAddress = ":" + env
	}
}
