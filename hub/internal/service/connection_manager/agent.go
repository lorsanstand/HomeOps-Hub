package connection_manager

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/shared/domain"
	"github.com/rs/zerolog"
)

type AgentConnection struct {
	stream             StreamConn
	heartbeat          heartbeatStore
	log                zerolog.Logger
	status             StatusAgent
	AgentID            string
	response           *ResponseStore
	ctx                context.Context
	cancel             context.CancelFunc
	heartbeatTimeoutMS int
}

func newAgentConnection(agentID string, stream StreamConn, heartbeat heartbeatStore, status StatusAgent, heartbeatTimeoutMS int, logger zerolog.Logger) *AgentConnection {
	response := NewResponseStore()
	logger = logger.With().Str("agentID", agentID).Logger()
	ctx, cancel := context.WithCancel(stream.Context())
	return &AgentConnection{stream: stream, response: response, heartbeat: heartbeat, log: logger, AgentID: agentID, status: status, ctx: ctx, cancel: cancel, heartbeatTimeoutMS: heartbeatTimeoutMS}
}

func (a *AgentConnection) Listen() error {
	heartbeatsCh := make(chan domainHub.CreateHeartbeatModel, 5)
	streamRecvCh := make(chan *pb.AgentEvent, 5)

	go a.listenHeartbeat(heartbeatsCh)
	go a.listenStream(streamRecvCh)

	defer func() {
		a.status.Offline()
		close(heartbeatsCh)
		a.Close()
	}()

	for {
		select {
		case <-a.ctx.Done():
			return a.ctx.Err()
		case msg, ok := <-streamRecvCh:
			if !ok {
				return nil
			}

			switch x := msg.Event.(type) {
			case *pb.AgentEvent_Heartbeat:
				heartbeat := toCreateHeartbeatModel(a.AgentID, x)
				heartbeatsCh <- heartbeat
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

func (a *AgentConnection) listenStream(ch chan *pb.AgentEvent) {
	defer close(ch)
	for {
		agentEvent, err := a.stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			a.log.Warn().Err(err).Msg("close stream")
			return
		}

		ch <- agentEvent
	}
}

func (a *AgentConnection) listenHeartbeat(heartbeats <-chan domainHub.CreateHeartbeatModel) {
	lastHeartbeat := 0
	timer := time.NewTicker(time.Duration(a.heartbeatTimeoutMS) * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if lastHeartbeat < 4 {
				lastHeartbeat += 1
				a.status.Offline()
				continue
			}

			a.log.Warn().Msg("agent not send heartbeat")
			a.Close()
			return
		case heartbeat, ok := <-heartbeats:
			if !ok {
				return
			}
			a.log.Debug().
				Float64("cpu usage", heartbeat.Metrics.CPUUsage).
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
		return domainHub.AgentResponse{}, ErrConnectionClose
	case <-ctx.Done():
		return domainHub.AgentResponse{}, ctx.Err()
	case response := <-ch:
		a.log.Info().Str("requestID", requestID).Msg("received response")
		return response, nil
	}
}

func (a *AgentConnection) Close() {
	a.cancel()
}
