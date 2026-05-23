package connection_manager

import (
	"sync"

	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
)

type AgentConnStore struct {
	mutex sync.RWMutex
	store map[string]*AgentConnection
}

func NewAgentConnStore() *AgentConnStore {
	return &AgentConnStore{store: make(map[string]*AgentConnection)}
}

func (a *AgentConnStore) Get(agentID string) *AgentConnection {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.store[agentID]
}

func (a *AgentConnStore) Add(agentID string, agentConn *AgentConnection) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.store[agentID] = agentConn
}

func (a *AgentConnStore) Delete(agentID string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.store, agentID)
}

func (a *AgentConnStore) Pop(agentID string) *AgentConnection {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	agent := a.store[agentID]
	delete(a.store, agentID)
	return agent
}

func (a *AgentConnStore) GetAllAgentID() []string {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var IDs []string
	for ID := range a.store {
		IDs = append(IDs, ID)
	}
	return IDs
}

type ResponseStore struct {
	store map[string]chan domainHub.AgentResponse
	mutex sync.RWMutex
}

func NewResponseStore() *ResponseStore {
	data := make(map[string]chan domainHub.AgentResponse)
	return &ResponseStore{store: data}
}

func (r *ResponseStore) Write(responseID string, channel chan domainHub.AgentResponse) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.store[responseID] = channel
}

func (r *ResponseStore) Read(responseID string) (chan domainHub.AgentResponse, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ch, ok := r.store[responseID]
	return ch, ok
}

func (r *ResponseStore) Delete(responseID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.store, responseID)
}
