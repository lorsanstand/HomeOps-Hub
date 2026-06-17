package command_store

import (
	"sync"

	"github.com/lorsanstand/HomeOps-Hub/agent/internal/domain"
)

type CommandStore struct {
	mutex sync.RWMutex
	store map[string]*domain.Command
}

func (c *CommandStore) Get(name string) (*domain.Command, bool) {
	c.mutex.RLock()
	function, ok := c.store[name]
	c.mutex.Unlock()
	return function, ok
}
