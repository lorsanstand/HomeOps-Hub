package connection_manager

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	"github.com/rs/zerolog"
	"gotest.tools/v3/assert"
)

type streamMock struct {
	recvCh  chan *pb.AgentEvent
	sendCh  chan *pb.ServerCommandRequest
	closeCh chan struct{}
	ctx     context.Context
}

func (f *streamMock) Context() context.Context {
	return f.ctx
}

func (f *streamMock) Send(request *pb.ServerCommandRequest) error {
	f.sendCh <- request
	return nil
}

func (f *streamMock) Recv() (*pb.AgentEvent, error) {
	select {
	case msg, ok := <-f.recvCh:
		if !ok {
			return nil, io.EOF
		}
		return msg, nil
	case <-f.ctx.Done():
		return nil, f.ctx.Err()
	}
}

func (f *streamMock) Close() error {
	select {
	case f.closeCh <- struct{}{}:
		close(f.recvCh)
	default:
	}
	return nil
}

type heartBeatMock struct {
	countUse int
	doneCh   chan struct{}
}

func (h *heartBeatMock) CreateHeartbeat(ctx context.Context, heartbeat domainHub.CreateHeartbeatModel) error {
	h.countUse += 1
	select {
	case h.doneCh <- struct{}{}:
	default:
	}
	return nil
}

type statusMock struct {
	online bool
	doneCh chan struct{}
}

func (s *statusMock) Offline() {
	s.online = false
}

func (s *statusMock) Online() {
	s.online = true
	select {
	case s.doneCh <- struct{}{}:
	default:
	}
}

func TestAgentConnection_Heartbeat(t *testing.T) {
	// Создаем вся поля для Agent Connection
	// Нужно как то вынести в отдельную функцию
	sendStream := make(chan *pb.ServerCommandRequest, 1)
	recvStream := make(chan *pb.AgentEvent)
	ctx, cancel := context.WithCancel(context.Background())

	stream := streamMock{recvCh: recvStream, sendCh: sendStream, ctx: ctx, closeCh: make(chan struct{}, 1)}
	heartbeat := heartBeatMock{doneCh: make(chan struct{}, 1)}
	status := statusMock{doneCh: make(chan struct{}, 1)}

	agent := newAgentConnection("123", &stream, &heartbeat, &status, 5000, zerolog.New(nil))
	go agent.Listen()

	recvStream <- &pb.AgentEvent{AgentId: "agent-1", Event: &pb.AgentEvent_Heartbeat{
		Heartbeat: &pb.Heartbeat{
			Timestamp: time.Now().Unix(),
			Metrics:   &pb.SystemMetrics{CpuUsage: 0.5, MemoryUsage: 0.3, DiskUsage: 0.7},
		}}}

	select {
	case <-heartbeat.doneCh:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for heartbeat")
	}

	select {
	case <-status.doneCh:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for status online")
	}

	assert.Equal(t, heartbeat.countUse, 1)
	assert.Equal(t, status.online, true)

	cancel()

	select {
	case <-stream.closeCh:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for close")
	}

	assert.Equal(t, status.online, false)
}

func TestAgentConnection_Execute(t *testing.T) {
	sendStream := make(chan *pb.ServerCommandRequest, 1)
	recvStream := make(chan *pb.AgentEvent)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	stream := streamMock{recvCh: recvStream, sendCh: sendStream, ctx: ctx}
	heartbeat := heartBeatMock{doneCh: make(chan struct{}, 1)}
	status := statusMock{doneCh: make(chan struct{}, 1)}

	agent := newAgentConnection("123", &stream, &heartbeat, &status, 5000, zerolog.New(nil))
	go agent.Listen()

	// Данные для проверки
	requestID := make(chan domainHub.AgentResponse)
	output := "test output"
	name := "test name"

	go func() {
		response, _ := agent.Execute(ctx, domainHub.AgentRequest{
			Name:    name,
			Args:    nil,
			TimeOut: 0,
		})

		requestID <- response
	}()

	request := <-sendStream
	assert.Equal(t, name, request.Name)

	recvStream <- &pb.AgentEvent{AgentId: "agent-1", Event: &pb.AgentEvent_CommandResponse{
		CommandResponse: &pb.CommandResponse{
			RequestId: request.RequestId,
			Success:   true,
			Output:    output,
		}}}

	select {
	case response := <-requestID:
		assert.Equal(t, output, response.Output)

	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for response")
	}
}

// Написать тест когда heartbeat не приходит и все закрывается
func TestAgentConnection_HeartbeatTimeout(t *testing.T) {
	sendStream := make(chan *pb.ServerCommandRequest, 1)
	recvStream := make(chan *pb.AgentEvent)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	stream := streamMock{recvCh: recvStream, sendCh: sendStream, ctx: ctx, closeCh: make(chan struct{}, 1)}
	heartbeat := heartBeatMock{doneCh: make(chan struct{}, 1)}
	status := statusMock{doneCh: make(chan struct{}, 1)}
	var wg sync.WaitGroup

	agent := newAgentConnection("123", &stream, &heartbeat, &status, 200, zerolog.New(nil))

	wg.Add(2)
	go func() {
		err := agent.Listen()
		assert.NilError(t, err)
		wg.Done()
	}()

	go func() {
		_, err := agent.Execute(ctx, domainHub.AgentRequest{
			Name:    "test",
			Args:    nil,
			TimeOut: 0,
		})
		assert.ErrorIs(t, err, ConnectionCloseErr)
		wg.Done()
	}()

	wg.Wait()

	select {
	case <-stream.closeCh:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for close")
	}
}

//Написать тест при закрытии соединения Execute завершается
