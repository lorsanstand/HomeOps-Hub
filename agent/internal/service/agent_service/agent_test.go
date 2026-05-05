package agent_service

import (
	"context"
	"errors"
	"testing"

	"github.com/lorsanstand/HomeOps-Hub/agent/internal/utils/config_yaml"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
	"github.com/rs/zerolog"
)

type CollectorMock struct {
	host domain.HostInfo
	caps []domain.Capability
}

func (c *CollectorMock) GatherInfoSystem() (domain.HostInfo, []domain.Capability) {
	return c.host, c.caps
}

type ConnectionMock struct {
	regAgentErr error
	regResp     domain.RegisterAgentResponse
	regData     domain.RegisterAgentRequest
}

func (c *ConnectionMock) RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error) {
	c.regData = RegisterData
	return c.regResp, c.regAgentErr
}

type SettingsMock struct {
	insertErr error
	agentID   string
	countUse  int
}

func (s *SettingsMock) InsertAgentID(agentID string) error {
	s.agentID = agentID
	s.countUse++
	return s.insertErr
}

func (s *SettingsMock) GetAgentID() string {
	return s.agentID
}

func TestAgentService_RegisterAgentConn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		wantErr          error
		insertAgentIDUse int
		settings         SettingsMock
		collector        CollectorMock
		conn             ConnectionMock
		cfg              config_yaml.AgentConfig
	}{
		{
			name:             "success",
			wantErr:          nil,
			insertAgentIDUse: 1,
			settings:         SettingsMock{agentID: "", insertErr: nil},
			collector: CollectorMock{
				host: domain.HostInfo{System: "Linux", Hostname: "test", Arch: "x64"},
				caps: []domain.Capability{
					{Available: true, Version: "0", Name: "testCaps", Reason: ""},
				},
			},
			conn: ConnectionMock{regAgentErr: nil, regResp: domain.RegisterAgentResponse{AgentID: "123", Heartbeat: 4}},
			cfg:  config_yaml.AgentConfig{AppName: "test"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			svc := NewAgentService(&tt.collector,
				&tt.conn,
				&tt.settings,
				&tt.cfg,
				zerolog.New(nil),
			)

			err := svc.RegisterAgentConn(ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got: %v", tt.wantErr, err)
			}

			if tt.insertAgentIDUse != tt.settings.countUse {
				t.Errorf("expected count insert agent id %v, got: %v", tt.insertAgentIDUse, tt.settings.countUse)
			}

			if tt.settings.agentID != tt.conn.regResp.AgentID {
				t.Errorf("expected insert agent id %v, got: %v", tt.conn.regResp.AgentID, tt.settings.agentID)
			}

			if tt.cfg.AppName != tt.conn.regData.AgentName {
				t.Fatalf("expected agent name %v, got: %v", tt.cfg.AppName, tt.conn.regData.AgentName)
			}

		})
	}
}
