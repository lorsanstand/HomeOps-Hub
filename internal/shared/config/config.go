package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBPassword string `env:"DB_PASS"`
	DBUser     string `env:"DB_USER"`
	DBName     string `env:"DB_NAME"`
	LogLevel   string `env:"LOG_LEVEL" env-default:"INFO"`
	Mode       string `env:"MODE" env-default:"DEV"`
	Port       int    `env:"PORT" env-default:"9000"`
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

func (c *Config) GetURLPostgres() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName)
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
