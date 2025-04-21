package config

import (
	"os"
	"strconv"
)

type Config struct {
	Http     *Http
	Database *Database
	Env      string
}

type Http struct {
	Addr string
}

type Database struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func Load() *Config {
	http := &Http{
		Addr: getString("HTTP_ADDR", ":8080"),
	}

	database := &Database{
		Addr:         getString("DATABASE_ADDR", "postgres://admin:adminpassword@localhost:5432/ecom?sslmode=disable"),
		MaxOpenConns: getInt("DATABASE_MAX_OPEN_CONNS", 30),
		MaxIdleConns: getInt("DATABASE_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  getString("DATABASE_MAX_IDLE_TIME", "15min"),
	}

	return &Config{
		Http:     http,
		Database: database,
		Env:      getString("ENV", "development"),
	}
}

func getString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
