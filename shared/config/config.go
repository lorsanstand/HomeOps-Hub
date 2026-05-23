package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" env-default:"INFO"`
	Mode     string `env:"MODE" env-default:"DEV"`
	Port     int    `env:"PORT" env-default:"9000"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Printf("failed read config: %v", err)

		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("env read failed: %v", err)
		}
	}

	return &cfg, nil
}

func (c *Config) GetLogLevel() zerolog.Level {
	level, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		return zerolog.InfoLevel
	}
	return level
}

func (c *Config) GetMode() string {
	return c.Mode
}
