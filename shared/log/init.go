package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type cfgLogStore interface {
	GetLogLevel() zerolog.Level
	GetMode() string
}

func NewLogger(cfg cfgLogStore) zerolog.Logger {
	var output io.Writer = os.Stdout

	if cfg.GetMode() != "PROD" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Kitchen,
		}
	}

	level := cfg.GetLogLevel()

	return zerolog.New(output).Level(level).With().Timestamp().Logger()
}
