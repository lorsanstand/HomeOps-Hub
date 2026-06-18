package domain

type RegisterAgentRequest struct {
	AgentID      string
	AgentName    string
	AgentVersion string
	Host         HostInfo
	Capabilities map[string]Capability
}

type HostInfo struct {
	System   string
	Hostname string
	Arch     string
}

type Capability struct {
	Available bool
	Version   string
	Reason    string
	Command   map[string]CapabilityCommand
}

type CapabilityCommand struct {
	OptionalArgs map[string]CommandArgs
	RequiredArgs map[string]CommandArgs
	Version      string
	Description  string
	TypeOutput   string
}

type CommandArgs struct {
	Type        string
	Description string
	Default     string
	Enum        []string
	Validation  ArgValidation
}

type ArgValidation struct {
	MinValue    int
	MaxValue    int
	Pattern     string
	AllowedExts []string
}

type RegisterAgentResponse struct {
	Heartbeat int
	AgentID   string
}
