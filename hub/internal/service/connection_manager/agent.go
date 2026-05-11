package connection_manager

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/service/connection_manager/store"
	"github.com/rs/zerolog"
)

type statusAgent interface {
	Offline()
	Online()
}

// использовать sync.Pool что бы переиспользвоать этот обьект
type AgentConnection struct {
	stream    streamConn
	heartbeat heartbeatStore
	log       zerolog.Logger
	status    statusAgent
	AgentID   string
	response  *store.ResponseStore
}

func newAgentConnection(agentID string, stream streamConn, heartbeat heartbeatStore, status statusAgent, logger zerolog.Logger) *AgentConnection {
	response := store.NewResponseStore()
	return &AgentConnection{stream: stream, response: response, heartbeat: heartbeat, log: logger, AgentID: agentID, status: status}
}

func (a *AgentConnection) Listen() error {
	ctx := a.stream.Context()
	defer a.status.Offline()

	heartbeatsChan := make(chan domainHub.CreateHeartbeatModel, 5)
	go a.listenHeartbeat(ctx, heartbeatsChan)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			agentEvent, err := a.stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("stream error: %w", err)
			}

			switch x := agentEvent.Event.(type) {
			case *pb.AgentEvent_Heartbeat:
				heartbeat := toCreateHeartbeatModel(a.AgentID, x)

				a.log.Debug().
					Str("agentID", heartbeat.AgentID).
					Float64("cpu usage", heartbeat.Metrics.CpuUsage).
					Float64("disk usage", heartbeat.Metrics.DiskUsage).
					Float64("memory usage", heartbeat.Metrics.MemoryUsage).Msg("")

				heartbeatsChan <- heartbeat
			}
		}
	}
}

func (a *AgentConnection) listenHeartbeat(ctx context.Context, heartbeats <-chan domainHub.CreateHeartbeatModel) {
	lastHeartbeat := 0
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if lastHeartbeat < 30 {
				lastHeartbeat += 5
				a.status.Offline()
				continue
			}
			a.log.Warn().Str("agentID", a.AgentID).Msg("agent did not send heartbeat")
			a.stream.Close()
			return
		case heartbeat := <-heartbeats:
			a.status.Online()
			lastHeartbeat = 0
			err := a.heartbeat.CreateHeartbeat(ctx, heartbeat)
			if err != nil {
				a.log.Error().Err(err).Str("agentID", heartbeat.AgentID).Msg("failed to write heartbeat")
			}
		case <-ctx.Done():
			return
		}
	}
}
