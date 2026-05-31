package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	return Config{
		Port: getENV("PORT", ":8080"),
	}
}

func getENV(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}

	return env
}
