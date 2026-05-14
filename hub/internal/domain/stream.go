package domain

import "time"

type AgentRequest struct {
	Name    string
	Args    map[string]string
	TimeOut int
}

type AgentResponse struct {
	Success    bool
	Output     string
	Error      string
	ExecTimeMS int
}

type AgentAlert struct {
	Timestamp   int
	Level       string
	Title       string
	Description string
}

type SystemMetrics struct {
	CpuUsage    float64
	MemoryUsage float64
	DiskUsage   float64
}

type HeartbeatModel struct {
	ID        int
	AgentID   string
	Timestamp time.Time
	Metrics   SystemMetrics
}

type CreateHeartbeatModel struct {
	AgentID   string
	Timestamp time.Time
	Metrics   SystemMetrics
}
