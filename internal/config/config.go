package config

import (
	"os"
	"time"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() Config {
	readTimeout, _ := time.ParseDuration(getENV("READ_TIMEOUT", "5s"))
	writeTimeout, _ := time.ParseDuration(getENV("WRITE_TIMEOUT", "10s"))
	idleTimeout, _ := time.ParseDuration(getENV("IDLE_TIMEOUT", "60s"))

	return Config{
		Port:         getENV("PORT", ":8080"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func getENV(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}

	return env
}
