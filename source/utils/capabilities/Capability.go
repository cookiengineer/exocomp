package capabilities

import "exocomp/types"

// Capability is the wire type owned by exocomp/types. The alias keeps the
// capability package API concise without duplicating the type.
type Capability = types.Capability

const (
	CapabilityFilesRead      = types.CapabilityFilesRead
	CapabilityFilesWrite     = types.CapabilityFilesWrite
	CapabilityNetworkConnect = types.CapabilityNetworkConnect
)
