// Package config
package config

import (
	"os"
)

type Config struct {
	DBPassword string
	DBPort string 
}

func Load() *Config {
	cfg := &Config{
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBPort: os.Getenv("DB_PORT"),
	}

	if cfg.DBPassword == "" {
		return nil
	}

	if cfg.DBPort == "" {
		return nil
	}

	return cfg
}
