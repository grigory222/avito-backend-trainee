package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		HTTP HTTP
		Log  Log
		DB   DB
	}

	HTTP struct {
		Host string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
		Port int    `env:"HTTP_PORT" envDefault:"8080"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	DB struct {
		Host     string `env:"DB_HOST,required"`
		Port     int    `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Name     string `env:"DB_NAME,required"`
		Password string `env:"DB_PASSWORD,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
