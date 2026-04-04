package app

import (
	"github.com/lorsanstand/HomeOps-Hub/internal/shared/config"
	"github.com/rs/zerolog"
)

type App struct {
	cfg *config.Config
	log *zerolog.Logger
}
