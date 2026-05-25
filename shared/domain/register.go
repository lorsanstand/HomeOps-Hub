package domain

type RegisterAgentRequest struct {
	AgentID      string
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
	Command   []CapabilityCommand
}

type CapabilityCommand struct {
	Name         string
	OptionalArgs []CommandArgs
	RequiredArgs []CommandArgs
	Version      string
	Description  string
	TypeOutput   string
}

type CommandArgs struct {
	Name        string
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
