package tools

import "exocomp/types"
import "fmt"
import "strings"

// permissionSummary formats the granted and denied capabilities of a permission
// request.
func permissionSummary(granted []types.Capability, denied []types.Capability) string {

	lines := make([]string, 0)

	if len(granted) > 0 {
		noun := "capabilities"
		if len(granted) == 1 {
			noun = "capability"
		}
		lines = append(lines, fmt.Sprintf("humans.AskPermission: Granted %d %s.", len(granted), noun))
		has_network := false
		for _, capability := range granted {
			lines = append(lines, "- "+capability.Label())
			if capability.Type == types.CapabilityNetworkConnect {
				has_network = true
			}
		}
		if has_network == true {
			lines = append(lines, "NOTE: network.Connect enables outbound network access for every sandboxed program; the scoped host is enforced for websites.* only.")
		}
	} else {
		lines = append(lines, "humans.AskPermission: No capabilities granted.")
	}

	if len(denied) > 0 {
		lines = append(lines, "Denied:")
		for _, capability := range denied {
			lines = append(lines, "- "+capability.Label())
		}
	}

	return strings.Join(lines, "\n")

}
