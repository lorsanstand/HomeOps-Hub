package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type cfgLogStore interface {
	GetLogLevel() zerolog.Level
	GetMode() string
}

func Init(cfg cfgLogStore) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if cfg.GetMode() != "PROD" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Kitchen,
		})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
