package connection_manager

import (
	"context"
	"errors"
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
	mu      sync.Mutex
	sendErr error
	recvErr error
	closeOnce sync.Once
}

func (f *streamMock) Context() context.Context {
	return f.ctx
}

func (f *streamMock) Send(request *pb.ServerCommandRequest) error {
	f.mu.Lock()
	err := f.sendErr
	f.mu.Unlock()
	if err != nil {
		return err
	}
	f.sendCh <- request
	return nil
}

func (f *streamMock) Recv() (*pb.AgentEvent, error) {
	f.mu.Lock()
	recvErr := f.recvErr
	f.mu.Unlock()
	if recvErr != nil {
		return nil, recvErr
	}
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
	default:
	}
	f.closeOnce.Do(func() {
		close(f.recvCh)
	})
	return nil
}

type heartBeatMock struct {
	mu       sync.Mutex
	countUse int
	doneCh   chan struct{}
	err      error
}

func (h *heartBeatMock) CreateHeartbeat(ctx context.Context, heartbeat domainHub.CreateHeartbeatModel) error {
	h.mu.Lock()
	h.countUse += 1
	err := h.err
	h.mu.Unlock()
	select {
	case h.doneCh <- struct{}{}:
	default:
	}
	return err
}

type statusMock struct {
	mu     sync.Mutex
	online bool
	doneCh chan struct{}
}

func (s *statusMock) Offline() {
	s.mu.Lock()
	s.online = false
	s.mu.Unlock()
}

func (s *statusMock) Online() {
	s.mu.Lock()
	s.online = true
	s.mu.Unlock()
	select {
	case s.doneCh <- struct{}{}:
	default:
	}
}

func (s *statusMock) IsOnline() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.online
}

type agentTestHarness struct {
	ctx       context.Context
	cancel    context.CancelFunc
	stream    *streamMock
	heartbeat *heartBeatMock
	status    *statusMock
	agent     *AgentConnection
	recvCh    chan *pb.AgentEvent
	sendCh    chan *pb.ServerCommandRequest
}

func newAgentTestHarness(t *testing.T, heartbeatTimeoutMS int) *agentTestHarness {
	t.Helper()
	sendStream := make(chan *pb.ServerCommandRequest, 4)
	recvStream := make(chan *pb.AgentEvent, 4)
	ctx, cancel := context.WithCancel(context.Background())

	stream := &streamMock{recvCh: recvStream, sendCh: sendStream, ctx: ctx, closeCh: make(chan struct{}, 1)}
	heartbeat := &heartBeatMock{doneCh: make(chan struct{}, 2)}
	status := &statusMock{doneCh: make(chan struct{}, 2)}

	agent := newAgentConnection("123", stream, heartbeat, status, heartbeatTimeoutMS, zerolog.New(nil))

	t.Cleanup(func() {
		cancel()
	})

	return &agentTestHarness{
		ctx: ctx, cancel: cancel, stream: stream, heartbeat: heartbeat, status: status,
		agent: agent, recvCh: recvStream, sendCh: sendStream,
	}
}

func waitFor(t *testing.T, ch <-chan struct{}, timeout time.Duration, message string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(timeout):
		t.Fatal(message)
	}
}

func waitForClose(t *testing.T, closeCh <-chan struct{}, timeout time.Duration) {
	t.Helper()
	select {
	case <-closeCh:
	case <-time.After(timeout):
		t.Fatal("timeout waiting for close")
	}
}

func commandResponseEvent(requestID, output string) *pb.AgentEvent {
	return &pb.AgentEvent{
		AgentId: "agent-1",
		Event: &pb.AgentEvent_CommandResponse{
			CommandResponse: &pb.CommandResponse{
				RequestId: requestID,
				Success:   true,
				Output:    output,
			},
		},
	}
}

func TestAgentConnection_Heartbeat(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	go h.agent.Listen()

	h.recvCh <- &pb.AgentEvent{AgentId: "agent-1", Event: &pb.AgentEvent_Heartbeat{
		Heartbeat: &pb.Heartbeat{
			Timestamp: time.Now().Unix(),
			Metrics:   &pb.SystemMetrics{CpuUsage: 0.5, MemoryUsage: 0.3, DiskUsage: 0.7},
		}}}

	waitFor(t, h.heartbeat.doneCh, 500*time.Millisecond, "timeout waiting for heartbeat")
	waitFor(t, h.status.doneCh, 500*time.Millisecond, "timeout waiting for status online")

	h.heartbeat.mu.Lock()
	count := h.heartbeat.countUse
	h.heartbeat.mu.Unlock()

	assert.Equal(t, count, 1)
	assert.Equal(t, h.status.IsOnline(), true)

	h.cancel()
	waitForClose(t, h.stream.closeCh, 500*time.Millisecond)
	assert.Equal(t, h.status.IsOnline(), false)
}

func TestAgentConnection_Execute(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	go h.agent.Listen()

	// Данные для проверки
	requestID := make(chan domainHub.AgentResponse)
	output := "test output"
	name := "test name"

	go func() {
		response, _ := h.agent.Execute(h.ctx, domainHub.AgentRequest{
			Name:    name,
			Args:    nil,
			TimeOut: 0,
		})

		requestID <- response
	}()

	request := <-h.sendCh
	assert.Equal(t, name, request.Name)

	h.recvCh <- &pb.AgentEvent{AgentId: "agent-1", Event: &pb.AgentEvent_CommandResponse{
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

func TestAgentConnection_HeartbeatTimeout(t *testing.T) {
	h := newAgentTestHarness(t, 200)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		err := h.agent.Listen()
		assert.NilError(t, err)
		wg.Done()
	}()

	go func() {
		_, err := h.agent.Execute(h.ctx, domainHub.AgentRequest{
			Name:    "test",
			Args:    nil,
			TimeOut: 0,
		})
		assert.ErrorIs(t, err, ConnectionCloseErr)
		wg.Done()
	}()

	wg.Wait()

	waitForClose(t, h.stream.closeCh, 500*time.Millisecond)
}

func TestAgentConnection_ConnectionClose(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		err := h.agent.Listen()
		assert.ErrorIs(t, err, context.Canceled)
		wg.Done()
	}()

	go func() {
		_, err := h.agent.Execute(context.Background(), domainHub.AgentRequest{
			Name:    "test",
			Args:    nil,
			TimeOut: 0,
		})
		assert.ErrorIs(t, err, ConnectionCloseErr)
		wg.Done()
	}()

	h.cancel()

	wg.Wait()

	waitForClose(t, h.stream.closeCh, 500*time.Millisecond)
}

func TestAgentConnection_ExecuteClose(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	ctxExecute, cancelExecute := context.WithCancel(context.Background())
	t.Cleanup(cancelExecute)

	executeCh := make(chan struct{})
	go h.agent.Listen()

	go func() {
		_, err := h.agent.Execute(ctxExecute, domainHub.AgentRequest{
			Name:    "test",
			Args:    nil,
			TimeOut: 0,
		})

		assert.ErrorIs(t, err, context.Canceled)
		executeCh <- struct{}{}
	}()

	cancelExecute()
	waitFor(t, executeCh, 500*time.Millisecond, "timeout waiting for execute close")
}

func TestAgentConnection_ListenEOF(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	h.stream.Close()

	err := h.agent.Listen()
	assert.NilError(t, err)
	waitForClose(t, h.stream.closeCh, 500*time.Millisecond)
}

func TestAgentConnection_ListenRecvError(t *testing.T) {
	h := newAgentTestHarness(t, 5000)

	recvErr := errors.New("recv failure")
	h.stream.mu.Lock()
	h.stream.recvErr = recvErr
	h.stream.mu.Unlock()

	err := h.agent.Listen()
	assert.ErrorIs(t, err, recvErr)
}

func TestAgentConnection_ExecuteSendError(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	h.stream.mu.Lock()
	h.stream.sendErr = errors.New("send failure")
	h.stream.mu.Unlock()

	_, err := h.agent.Execute(h.ctx, domainHub.AgentRequest{Name: "test"})
	assert.ErrorContains(t, err, "execute command")
}

func TestAgentConnection_ExecuteContextCanceled(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := h.agent.Execute(ctx, domainHub.AgentRequest{Name: "test"})
	assert.ErrorIs(t, err, context.Canceled)
}

func TestAgentConnection_ExecuteConnectionCanceled(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	h.cancel()

	_, err := h.agent.Execute(context.Background(), domainHub.AgentRequest{Name: "test"})
	assert.ErrorIs(t, err, ConnectionCloseErr)
}

func TestAgentConnection_UnknownResponseID(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	go h.agent.Listen()

	h.recvCh <- &pb.AgentEvent{AgentId: "agent-1", Event: &pb.AgentEvent_CommandResponse{
		CommandResponse: &pb.CommandResponse{
			RequestId: "unknown",
			Success:   true,
			Output:    "ok",
		}}}

	h.cancel()
	waitForClose(t, h.stream.closeCh, 500*time.Millisecond)
}

func TestAgentConnection_HeartbeatErrorDoesNotStop(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	h.heartbeat.mu.Lock()
	h.heartbeat.err = errors.New("db error")
	h.heartbeat.mu.Unlock()

	go h.agent.Listen()
	h.recvCh <- &pb.AgentEvent{AgentId: "agent-1", Event: &pb.AgentEvent_Heartbeat{
		Heartbeat: &pb.Heartbeat{
			Timestamp: time.Now().Unix(),
			Metrics:   &pb.SystemMetrics{CpuUsage: 0.2, MemoryUsage: 0.1, DiskUsage: 0.3},
		}}}

	waitFor(t, h.heartbeat.doneCh, 500*time.Millisecond, "timeout waiting for heartbeat")
	h.cancel()
}

func TestAgentConnection_ConcurrentExecute(t *testing.T) {
	h := newAgentTestHarness(t, 5000)
	go h.agent.Listen()

	responses := make(chan domainHub.AgentResponse, 2)

	go func() {
		resp, _ := h.agent.Execute(h.ctx, domainHub.AgentRequest{Name: "cmd-1"})
		responses <- resp
	}()
	go func() {
		resp, _ := h.agent.Execute(h.ctx, domainHub.AgentRequest{Name: "cmd-2"})
		responses <- resp
	}()

	first := <-h.sendCh
	second := <-h.sendCh

	// ответы приходят в обратном порядке
	h.recvCh <- commandResponseEvent(second.RequestId, "second")
	h.recvCh <- commandResponseEvent(first.RequestId, "first")

	resp1 := <-responses
	resp2 := <-responses

	assert.Assert(t, resp1.Output == "first" || resp1.Output == "second")
	assert.Assert(t, resp2.Output == "first" || resp2.Output == "second")
	assert.Assert(t, resp1.Output != resp2.Output)

	h.cancel()
}
