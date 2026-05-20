package connection_manager

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
)

const heartbeatTimeoutMS = 6000

type ConnectionManager struct {
	heartbeat      heartbeatStore
	log            zerolog.Logger
	status         statusNotifier
	agentConnStore *AgentConnStore
}

func NewConnectionManager(heartbeat heartbeatStore, status statusNotifier, logger zerolog.Logger) *ConnectionManager {
	return &ConnectionManager{heartbeat: heartbeat, log: logger, status: status, agentConnStore: NewAgentConnStore()}
}

func (c *ConnectionManager) NewConnection(stream streamConn) error {
	AgentID, err := agentIDFromMetadata(stream.Context())
	if err != nil {
		c.log.Error().Err(err).Msg("missing agent id in metadata")
		return fmt.Errorf("get agent id: %w", err)
	}
	c.log.Info().Str("agentID", AgentID).Msg("connection accepted")

	status := c.status.New(AgentID)

	agent := newAgentConnection(AgentID, stream, c.heartbeat, status, heartbeatTimeoutMS, c.log)
	c.agentConnStore.Add(AgentID, agent)
	go func() {
		c.log.Debug().Str("agentID", AgentID).Msg("start listening")
		err := agent.Listen()
		if err != nil {
			c.log.Error().Err(err).Msg("listening agent stopped")
		}
		c.agentConnStore.Delete(AgentID)
	}()

	return nil
}

func (c *ConnectionManager) GetConnection(AgentID string) (*AgentConnection, error) {
	agent := c.agentConnStore.Get(AgentID)
	if agent == nil {
		return nil, ErrNotFoundConn
	}

	return agent, nil
}

func (c *ConnectionManager) GetAllAgentID() []string {
	return c.agentConnStore.GetAllAgentID()
}

func agentIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not found")
	}

	values := md.Get("agent-id")
	if len(values) == 0 || values[0] == "" {
		return "", fmt.Errorf("agent-id not found")
	}

	return values[0], nil
}
