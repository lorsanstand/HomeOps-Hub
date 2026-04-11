package config

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestConfig_GetLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cfg          Config
		wantLogLevel zerolog.Level
	}{
		{
			name:         "success",
			cfg:          Config{LogLevel: "DEBUG"},
			wantLogLevel: zerolog.DebugLevel,
		},
		{
			name:         "failed parse",
			cfg:          Config{LogLevel: "TEST"},
			wantLogLevel: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logLevel := tt.cfg.GetLogLevel()

			if logLevel != tt.wantLogLevel {
				t.Fatalf("expected %v, got: %v", tt.wantLogLevel, logLevel)
			}
		})
	}
}

func TestConfig_GetURLPostgres(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		cfg             Config
		wantURLPostgres string
	}{
		{
			name: "success",
			cfg: Config{
				DBHost:     "TestHost",
				DBName:     "TestName",
				DBPassword: "TestPassword",
				DBPort:     1234,
				DBUser:     "TestUser",
			},
			wantURLPostgres: "postgres://TestUser:TestPassword@TestHost:1234/TestName?sslmode=disable",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			url := tt.cfg.GetURLPostgres()

			if url != tt.wantURLPostgres {
				t.Fatalf("expected %v, got: %v", tt.wantURLPostgres, url)
			}
		})
	}
}
