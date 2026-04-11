package hub_service

type AgentRegistrationData struct {
	AgentID   string
	AgentName string
	Config    AgentConfig
}

type AgentConfig struct {
	System   string
	Docker   bool
	Hostname string
	Version  string
	Arch     string
}
