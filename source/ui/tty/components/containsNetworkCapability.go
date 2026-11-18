package components

import "exocomp/types"

func containsNetworkCapability(capabilities []types.Capability) bool {

	for _, capability := range capabilities {

		if capability.Type == types.CapabilityNetworkConnect {
			return true
		}

	}

	return false

}
