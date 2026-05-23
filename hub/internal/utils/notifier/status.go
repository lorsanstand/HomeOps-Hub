package notifier

import "github.com/lorsanstand/HomeOps-Hub/hub/internal/service/connection_manager"

// Временная заглушка
type StatusNotifier struct {
}

type Status struct {
	agentID string
	online  bool
}

func NewStatusNotifier() *StatusNotifier {
	return &StatusNotifier{}
}

func (s *StatusNotifier) New(agentID string) connection_manager.StatusAgent {
	return &Status{agentID: agentID}
}

func (s *Status) Online() {
	s.online = true
}

func (s *Status) Offline() {
	s.online = false
}
