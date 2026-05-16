package store

import (
	"sync"

	"github.com/lorsanstand/HomeOps-Hub/hub/internal/service/connection_manager"
)

type AgentConnStore struct {
	mutex sync.RWMutex
	store map[string]*connection_manager.AgentConnection
}

func NewAgentConnStore() *AgentConnStore {
	return &AgentConnStore{store: make(map[string]*connection_manager.AgentConnection)}
}

func (a *AgentConnStore) Get(agentID string) *connection_manager.AgentConnection {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.store[agentID]
}

func (a *AgentConnStore) Add(agentConn *connection_manager.AgentConnection) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.store[agentConn.AgentID] = agentConn
}

func (a *AgentConnStore) Delete(agentID string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.store, agentID)
}

func (a *AgentConnStore) Pop(agentID string) *connection_manager.AgentConnection {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	agent := a.store[agentID]
	delete(a.store, agentID)
	return agent
}
