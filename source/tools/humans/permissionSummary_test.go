package humans

import "exocomp/types"
import "strings"
import "testing"

func TestPermissionSummary_NetworkNote(t *testing.T) {

	result := permissionSummary([]types.Capability{
		{Type: types.CapabilityNetworkConnect, Scope: "example.com"},
	}, []types.Capability{})

	if strings.Contains(result, "network.Connect example.com") == false {
		t.Errorf("Expected network capability in %q", result)
	}

	if strings.Contains(result, "outbound network access") == false {
		t.Errorf("Expected network note in %q", result)
	}

	result2 := permissionSummary([]types.Capability{
		{Type: types.CapabilityFilesRead, Scope: "/tmp/example"},
	}, []types.Capability{})

	if strings.Contains(result2, "outbound network access") == true {
		t.Errorf("Expected no network note in %q", result2)
	}

}
