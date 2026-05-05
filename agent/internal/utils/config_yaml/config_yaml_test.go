package config_yaml

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestAgentConfig_GetLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cfg          AgentConfig
		wantLogLevel zerolog.Level
	}{
		{
			name:         "success",
			cfg:          AgentConfig{LogLevel: "DEBUG"},
			wantLogLevel: zerolog.DebugLevel,
		},
		{
			name:         "failed parse",
			cfg:          AgentConfig{LogLevel: "TEST"},
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
