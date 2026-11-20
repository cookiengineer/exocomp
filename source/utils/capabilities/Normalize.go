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

		if scope != "" {

			if filepath.IsAbs(scope) == true {

				cleaned := filepath.Clean(scope)

				if IsPathDenied(cleaned) == true {
					return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Path is not allowed.", scope)
				}

				resolved, err0 := filepath.EvalSymlinks(cleaned)

				if err0 != nil {
					return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Path does not exist.", scope)
				}

				stat, err1 := os.Stat(resolved)

				if err1 != nil {
					return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Path does not exist.", scope)
				}

				if stat.IsDir() == false && stat.Mode().IsRegular() == false {
					return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Path must be a file or folder.", scope)
				}

				if IsPathDenied(resolved) == true {
					return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Path is not allowed.", scope)
				}

				return Capability{
					Type:  capability.Type,
					Scope: resolved,
				}, nil

			} else {
				return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Path must be absolute.", scope)
			}

		} else {
			return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope: Path is empty.")
		}

	case CapabilityNetworkConnect:

		scope := strings.ToLower(strings.TrimSpace(capability.Scope))

		if scope != "" {

			if ip := net.ParseIP(scope); ip != nil {

				return Capability{
					Type:  capability.Type,
					Scope: ip.String(),
				}, nil

			} else if ip, network, err := net.ParseCIDR(scope); err == nil {

				_ = ip
				return Capability{
					Type:  capability.Type,
					Scope: network.String(),
				}, nil

			} else if strings.Contains(scope, "://") == false && strings.ContainsAny(scope, "/?# \t\r\n") == false {

				if IsValidDomain(scope) == true {
					return Capability{
						Type:  capability.Type,
						Scope: scope,
					}, nil
				} else {
					return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Use a domain, IP or CIDR only.", capability.Scope)
				}

			} else {
				return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope \"%s\": Use a domain, IP or CIDR only.", capability.Scope)
			}

		} else {
			return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid scope: Host is empty.")
		}

	default:
		return Capability{}, fmt.Errorf("capabilities.Normalize: Invalid capability type \"%s\".", capability.Type)

	}

}
