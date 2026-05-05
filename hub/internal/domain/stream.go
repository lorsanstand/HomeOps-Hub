package domain

type AgentRequest struct {
	RequestID string
	Name      string
	Args      map[string]string
	TimeOut   int
}

type AgentResponse struct {
	RequestID  string
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

type Heartbeat struct {
	Timestamp string
	Metrics   SystemMetrics
}
