package config

import (
	"os"
	"strconv"
)

var (
	HTTPAddress     = ":80" // Address where the server will run
	ShutdownTimeout = 5     // Time in seconds for correct server shutdown
)

func init() {
	if env := os.Getenv("ADDRESS"); env != "" {
		HTTPAddress = env
	} else if env := os.Getenv("PORT"); env != "" {
		HTTPAddress = ":" + env
	}

	if env := os.Getenv("SHUTDOWN_TIMEOUT"); env != "" {
		var err error
		ShutdownTimeout, err = strconv.Atoi(env)
		if err != nil {
			panic(err)
		}
	}
}
