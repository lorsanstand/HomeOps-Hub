package connection_manager

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
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
	ctx       context.Context
	cancel    context.CancelFunc
}

func newAgentConnection(agentID string, stream streamConn, heartbeat heartbeatStore, status statusAgent, logger zerolog.Logger) *AgentConnection {
	response := store.NewResponseStore()
	logger = logger.With().Str("agentID", agentID).Logger()
	ctx, cancel := context.WithCancel(stream.Context())
	return &AgentConnection{stream: stream, response: response, heartbeat: heartbeat, log: logger, AgentID: agentID, status: status, ctx: ctx, cancel: cancel}
}

func (a *AgentConnection) Listen() error {
	defer a.status.Offline()

	heartbeatsChan := make(chan domainHub.CreateHeartbeatModel, 5)
	go a.listenHeartbeat(heartbeatsChan)
	defer close(heartbeatsChan)

	for {
		select {
		case <-a.ctx.Done():
			err := a.stream.Close()
			if err != nil {
				a.log.Warn().Err(err).Msg("failed stream close")
			}
			return a.ctx.Err()
		default:
			agentEvent, err := a.stream.Recv()
			if err == io.EOF {
				a.cancel()
				return nil
			}
			if err != nil {
				a.cancel()
				return fmt.Errorf("stream: %w", err)
			}

			switch x := agentEvent.Event.(type) {
			case *pb.AgentEvent_Heartbeat:
				heartbeat := toCreateHeartbeatModel(a.AgentID, x)
				heartbeatsChan <- heartbeat
			case *pb.AgentEvent_CommandResponse:
				ch, ok := a.response.Read(x.CommandResponse.RequestId)
				if !ok {
					a.log.Warn().Str("requestID", x.CommandResponse.RequestId).Msg("not found channel for send response")
					continue
				}
				response := toAgentResponse(x)
				ch <- response
			}
		}
	}
}

func (a *AgentConnection) listenHeartbeat(heartbeats <-chan domainHub.CreateHeartbeatModel) {
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

			a.log.Warn().Msg("agent not send heartbeat")
			a.cancel()
			return
		case heartbeat := <-heartbeats:
			a.log.Debug().
				Float64("cpu usage", heartbeat.Metrics.CpuUsage).
				Float64("disk usage", heartbeat.Metrics.DiskUsage).
				Float64("memory usage", heartbeat.Metrics.MemoryUsage).Msg("")

			a.status.Online()
			lastHeartbeat = 0

			err := a.heartbeat.CreateHeartbeat(a.ctx, heartbeat)
			if err != nil {
				a.log.Error().Err(err).Msg("failed to write heartbeat")
			}
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *AgentConnection) Execute(ctx context.Context, request domainHub.AgentRequest) (domainHub.AgentResponse, error) {
	requestID := uuid.New().String()
	ch := make(chan domainHub.AgentResponse, 1)
	defer close(ch)

	a.response.Write(requestID, ch)
	defer a.response.Delete(requestID)

	err := a.stream.Send(new(toGRPCCommandRequest(requestID, request)))
	if err != nil {
		return domainHub.AgentResponse{}, fmt.Errorf("execute command: %w", err)
	}

	a.log.Info().Str("requestID", requestID).Str("command", request.Name).Msg("send command")

	select {
	case <-a.ctx.Done():
		return domainHub.AgentResponse{}, fmt.Errorf("connection close")
	case <-ctx.Done():
		return domainHub.AgentResponse{}, fmt.Errorf("request timeout")
	case response := <-ch:
		a.log.Info().Str("requestID", requestID).Msg("received response")
		return response, nil
	}
}
