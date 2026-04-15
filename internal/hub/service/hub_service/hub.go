package hub_service

import (
	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/utils/hasher"
	"github.com/rs/zerolog"
)

type HubService struct {
	log zerolog.Logger
}

func NewHubService(logger zerolog.Logger) *HubService {
	return &HubService{log: logger}
}

func (h *HubService) RegisterAgent(data domain.RegisterAgentRequest) domain.RegisterAgentResponse {
	AgentID := data.AgentId
	if data.AgentId == "" {
		var err error
		AgentID, err = hasher.MakeID(data.Host, data.AgentName)
		if err != nil {
			h.log.Error().Err(err).Msg("failed create agent id")
			AgentID = ""
		}
	}

	return domain.RegisterAgentResponse{AgentID: AgentID, Heartbeat: 5}
}
