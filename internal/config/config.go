package config

import (
	"os"
	"strconv"
)

var (
	HTTPAddress         = ":80"     // Address where the server will run
	ShutdownTimeout     = 5         // Time in seconds for correct server shutdown
	BotToken            string      // Secret authorization token
	CacheStorage        string      // The storage type used for caching images
	CacheMemoryCapacity = 100       // The memory storage capacity
	CacheFSPath         = "./cache" // The path to the directory used for file system-based cache storage.
)

func init() {
	var err error

	if env := os.Getenv("ADDRESS"); env != "" {
		HTTPAddress = env
	} else if env := os.Getenv("PORT"); env != "" {
		HTTPAddress = ":" + env
	}

	if env := os.Getenv("SHUTDOWN_TIMEOUT"); env != "" {
		ShutdownTimeout, err = strconv.Atoi(env)
		if err != nil {
			panic(err)
		}
	}

	BotToken = os.Getenv("BOT_TOKEN")

	CacheStorage = os.Getenv("CACHE_STORAGE")

	if env := os.Getenv("CACHE_MEMORY_CAPACITY"); env != "" {
		CacheMemoryCapacity, err = strconv.Atoi(env)
		if err != nil {
			panic(err)
		}
	}

	if env := os.Getenv("CACHE_FS_PATH"); env != "" {
		CacheFSPath = env
	}
}
