package config

import (
	"os"
)

type Config struct {
	DBUser, DBPass, DBHost, DBName string
	HTTPPort                       string
}

func Load() *Config {
	return &Config{
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASS"),
		DBHost:   os.Getenv("DB_HOST"),
		DBName:   os.Getenv("DB_NAME"),
		HTTPPort: os.Getenv("HTTP_PORT"),
	}
}
