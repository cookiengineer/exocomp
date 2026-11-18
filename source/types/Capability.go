package types

import "fmt"

const (
	CapabilityFilesRead      = "files.Read"
	CapabilityFilesWrite     = "files.Write"
	CapabilityNetworkConnect = "network.Connect"
)

type Capability struct {
	Type  string `json:"type"`
	Scope string `json:"scope"`
}

func (capability Capability) Key() string {
	return capability.Type + "\x00" + capability.Scope
}

func (capability Capability) Label() string {
	return fmt.Sprintf("%s %s", capability.Type, capability.Scope)
}
