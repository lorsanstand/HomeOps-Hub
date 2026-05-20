package connection_manager

import (
	"context"
	"io"
	"sync"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
)

type streamMock struct {
	recvCh    chan *pb.AgentEvent
	sendCh    chan *pb.ServerCommandRequest
	closeCh   chan struct{}
	ctx       context.Context
	mu        sync.Mutex
	sendErr   error
	recvErr   error
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

type statusNotifierMock struct {
	agentIDCh chan string
}

func (s *statusNotifierMock) New(AgentID string) statusAgent {
	select {
	case s.agentIDCh <- AgentID:
	default:

	}
	return &statusMock{doneCh: make(chan struct{}, 2)}
}
