package config

import "os"

type Config struct {
	Http *Http
	Env  string
}

type Http struct {
	Addr string
}

func Load() *Config {
	http := &Http{
		Addr: getString("HTTP_ADDR", ":8080"),
	}

	return &Config{
		Http: http,
		Env:  getString("ENV", "development"),
	}
}

func getString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
