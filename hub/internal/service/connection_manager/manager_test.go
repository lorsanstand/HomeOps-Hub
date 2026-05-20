package connection_manager

import (
	"context"
	"testing"
	"time"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"gotest.tools/v3/assert"
)

type connectionManagerTestHarness struct {
	heartbeat *heartBeatMock
	status    *statusNotifierMock
	manager   *ConnectionManager
}

func newConnectionManagerTestHarness(t *testing.T) *connectionManagerTestHarness {
	t.Helper()

	heartbeat := &heartBeatMock{doneCh: make(chan struct{}, 2)}
	status := &statusNotifierMock{agentIDCh: make(chan string, 1)}

	manager := NewConnectionManager(heartbeat, status, zerolog.New(nil))

	return &connectionManagerTestHarness{manager: manager, status: status, heartbeat: heartbeat}
}

func newMetadataAgentID(t *testing.T, agentID string) metadata.MD {
	t.Helper()

	return metadata.New(map[string]string{"agent-id": agentID})
}

func TestNewConnectionManager_NewConnection(t *testing.T) {
	h := newConnectionManagerTestHarness(t)
	agentID := "123"

	ctx := metadata.NewIncomingContext(context.Background(), newMetadataAgentID(t, agentID))

	stream := streamMock{ctx: ctx,
		recvCh:  make(chan *pb.AgentEvent, 1),
		sendCh:  make(chan *pb.ServerCommandRequest, 1),
		closeCh: make(chan struct{}, 1),
	}

	err := h.manager.NewConnection(&stream)
	assert.NilError(t, err)

	select {
	case ID := <-h.status.agentIDCh:
		require.Equal(t, agentID, ID)
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("get agent id for notifier")
	}

	agentIDs := h.manager.GetAllAgentID()
	assert.Equal(t, agentID, agentIDs[0])

	agent, err := h.manager.GetConnection(agentID)
	assert.NilError(t, err)
	require.NotNil(t, agent)

}

func TestNewConnectionManager_NewConnectionNotAgentID(t *testing.T) {
	h := newConnectionManagerTestHarness(t)

	stream := streamMock{ctx: context.Background(),
		recvCh:  make(chan *pb.AgentEvent, 1),
		sendCh:  make(chan *pb.ServerCommandRequest, 1),
		closeCh: make(chan struct{}, 1),
	}

	err := h.manager.NewConnection(&stream)
	assert.ErrorContains(t, err, "get agent id")
}

func TestNewConnectionManager_AgentNotFound(t *testing.T) {
	h := newConnectionManagerTestHarness(t)
	_, err := h.manager.GetConnection("123")
	assert.ErrorIs(t, ErrNotFoundConn, err)

	agentIDs := h.manager.GetAllAgentID()
	assert.Equal(t, len(agentIDs), 0)
}

func Test_agentIDFromMetadata(t *testing.T) {
	agentID := "123"

	ctx := metadata.NewIncomingContext(context.Background(), newMetadataAgentID(t, agentID))
	id, err := agentIDFromMetadata(ctx)
	assert.NilError(t, err)
	assert.Equal(t, id, agentID)
}

func Test_agentIDFromMetadata_MetadataNotFound(t *testing.T) {
	ctx := context.Background()
	_, err := agentIDFromMetadata(ctx)
	assert.ErrorContains(t, err, "metadata not found")
}

func Test_agentIDFromMetadata_AgentIDNotFound(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), newMetadataAgentID(t, ""))
	_, err := agentIDFromMetadata(ctx)
	assert.ErrorContains(t, err, "agent-id not found")
}
