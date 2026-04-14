package hub_service

import (
	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/rs/zerolog"
)

type HubService struct {
	log zerolog.Logger
}

func NewHubService(request domain.RegisterAgentRequest)
