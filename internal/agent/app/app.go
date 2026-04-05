package app

import (
	standartlog "log"

	"github.com/lorsanstand/HomeOps-Hub/internal/agent/rpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/config_yaml"
	log2 "github.com/lorsanstand/HomeOps-Hub/internal/shared/log"
	"github.com/rs/zerolog"
)

type App struct {
	log     zerolog.Logger
	cfg     *config_yaml.AgentConfig
	hubConn *rpc.Connection
}

func NewApp() *App {
	cfg, err := config_yaml.NewConfig()
	if err != nil {
		standartlog.Fatalf("failed get config: %v", err)
	}

	log := log2.NewLogger(cfg)

	return &App{cfg: cfg, log: log}
}

func (a *App) Run() {
	conn, err := rpc.NewConnectAgent(a.cfg.GetGRPCAddress())
	if err != nil {
		a.log.Error().Err(err)
		return
	}

	a.hubConn = conn
	r
}
