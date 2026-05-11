package store

import (
	"sync"

	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
)

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

func (r *ResponseStore) Read(responseID string) chan domainHub.AgentResponse {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.store[responseID]
}

func (r *ResponseStore) Delete(responseID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.store, responseID)
}
