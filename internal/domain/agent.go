package domain

type RegisterAgentRequest struct {
	AgentId      string
	AgentName    string
	AgentVersion string
	Host         HostInfo
	Capabilities []Capability
}

type HostInfo struct {
	System   string
	Hostname string
	Arch     string
}

type Capability struct {
	Available bool
	Version   string
	Name      string
	Reason    string
}

type RegisterAgentResponse struct {
	Heartbeat int
	AgentID   string
}
