package capabilities

import "fmt"
import "net"
import "os"
import "path/filepath"
import "strings"

// Normalize validates and canonicalizes one capability.
func Normalize(capability Capability) (Capability, error) {

	switch capability.Type {

	case CapabilityFilesRead, CapabilityFilesWrite:

		scope := strings.TrimSpace(capability.Scope)

		if scope == "" {
			return Capability{}, fmt.Errorf("Invalid scope: Path is empty.")
		}

		if filepath.IsAbs(scope) == false {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Path must be absolute.", scope)
		}

		cleaned := filepath.Clean(scope)

		if IsPathDenied(cleaned) == true {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Path is not allowed.", scope)
		}

		resolved, err0 := filepath.EvalSymlinks(cleaned)

		if err0 != nil {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Path does not exist.", scope)
		}

		stat, err1 := os.Stat(resolved)

		if err1 != nil {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Path does not exist.", scope)
		}

		if stat.IsDir() == false && stat.Mode().IsRegular() == false {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Path must be a file or folder.", scope)
		}

		if IsPathDenied(resolved) == true {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Path is not allowed.", scope)
		}

		return Capability{
			Type:  capability.Type,
			Scope: resolved,
		}, nil

	case CapabilityNetworkConnect:

		scope := strings.ToLower(strings.TrimSpace(capability.Scope))

		if scope == "" {
			return Capability{}, fmt.Errorf("Invalid scope: Host is empty.")
		}

		if ip := net.ParseIP(scope); ip != nil {
			return Capability{
				Type:  capability.Type,
				Scope: ip.String(),
			}, nil
		}

		if ip, network, err := net.ParseCIDR(scope); err == nil {
			_ = ip
			return Capability{
				Type:  capability.Type,
				Scope: network.String(),
			}, nil
		}

		if strings.Contains(scope, "://") == true || strings.ContainsAny(scope, "/?# \t\r\n") == true {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Use a domain, IP or CIDR only.", capability.Scope)
		}

		if IsValidDomain(scope) == true {
			return Capability{
				Type:  capability.Type,
				Scope: scope,
			}, nil
		} else {
			return Capability{}, fmt.Errorf("Invalid scope \"%s\": Use a domain, IP or CIDR only.", capability.Scope)
		}

	default:
		return Capability{}, fmt.Errorf("Invalid capability type \"%s\".", capability.Type)
	}

}
