package capabilities

import "net"
import "sort"
import "strings"
import "sync"

// Set stores granted capabilities for one agent session.
type Set struct {
	mutex        sync.RWMutex
	capabilities []Capability
}

// NewSet creates an empty capability set.
func NewSet() *Set {
	return &Set{
		capabilities: make([]Capability, 0),
	}
}

func (set *Set) Grant(capabilities []Capability) error {

	normalized := make([]Capability, 0, len(capabilities))

	for _, capability := range capabilities {

		tmp, err0 := Normalize(capability)

		if err0 == nil {
			normalized = append(normalized, tmp)
		} else {
			return err0
		}

	}

	set.mutex.Lock()
	defer set.mutex.Unlock()

	for _, capability := range normalized {

		found := false

		for _, existing := range set.capabilities {

			if existing.Key() == capability.Key() {
				found = true
				break
			}

		}

		if found == false {
			set.capabilities = append(set.capabilities, capability)
		}

	}

	sort.Slice(set.capabilities, func(a int, b int) bool {
		return set.capabilities[a].Key() < set.capabilities[b].Key()
	})

	return nil

}

func (set *Set) List() []Capability {

	if set == nil {
		return []Capability{}
	} else {

		set.mutex.RLock()
		defer set.mutex.RUnlock()

		result := make([]Capability, len(set.capabilities))
		copy(result, set.capabilities)

		return result

	}
}

func (set *Set) CanRead(path string) bool {

	if set != nil {

		canonical := CanonicalPath(path)

		set.mutex.RLock()
		defer set.mutex.RUnlock()

		for _, capability := range set.capabilities {

			if capability.Type == CapabilityFilesRead || capability.Type == CapabilityFilesWrite {

				if IsPathWithin(capability.Scope, canonical) == true {
					return true
				}

			}

		}

		return false

	} else {
		return false
	}
}

func (set *Set) CanWrite(path string) bool {

	if set != nil {

		canonical := CanonicalPath(path)

		set.mutex.RLock()
		defer set.mutex.RUnlock()

		for _, capability := range set.capabilities {

			if capability.Type == CapabilityFilesWrite {

				if IsPathWithin(capability.Scope, canonical) == true {
					return true
				}

			}

		}

		return false

	} else {
		return false
	}
}

// AllowsNetworkEgress reports whether any network.Connect capability is
// granted. The OS sandbox can only allow or deny network access as a whole, so
// a single grant enables egress for every sandboxed program of the session.
func (set *Set) AllowsNetworkEgress() bool {
	return set.HasNetwork()
}

func (set *Set) HasNetwork() bool {

	if set != nil {

		set.mutex.RLock()
		defer set.mutex.RUnlock()

		for _, capability := range set.capabilities {

			if capability.Type == CapabilityNetworkConnect {
				return true
			}

		}

		return false

	} else {
		return false
	}
}

func (set *Set) CanConnect(host string) bool {

	if set != nil {

		host = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(host), "."))

		if host != "" {

			host_ip := net.ParseIP(host)

			set.mutex.RLock()
			defer set.mutex.RUnlock()

			for _, capability := range set.capabilities {

				if capability.Type == CapabilityNetworkConnect {

					if host_ip != nil {

						if NetworkContains(capability.Scope, host_ip) == true {
							return true
						}

					} else if capability.Scope == host {
						return true
					}

				}

			}

			return false

		} else {
			return false
		}
	} else {
		return false
	}
}

func (set *Set) CanConnectIP(ip net.IP) bool {

	if set != nil && ip != nil {

		set.mutex.RLock()
		defer set.mutex.RUnlock()

		for _, capability := range set.capabilities {

			if capability.Type == CapabilityNetworkConnect {

				if NetworkContains(capability.Scope, ip) == true {
					return true
				}

			}

		}

		return false

	} else {
		return false
	}
}

func (set *Set) Covers(capability Capability) bool {

	if set != nil {

		normalized, err0 := Normalize(capability)

		if err0 == nil {

			set.mutex.RLock()
			defer set.mutex.RUnlock()

			for _, existing := range set.capabilities {

				if existing.Type == normalized.Type {

					if existing.Scope == normalized.Scope {
						return true
					} else if (existing.Type == CapabilityFilesRead || existing.Type == CapabilityFilesWrite) && IsPathWithin(existing.Scope, normalized.Scope) == true {
						return true
					}

				}

				if normalized.Type == CapabilityFilesRead && existing.Type == CapabilityFilesWrite {

					if IsPathWithin(existing.Scope, normalized.Scope) == true {
						return true
					}

				}

			}

			return false

		} else {
			return false
		}
	} else {
		return false
	}
}

func (set *Set) SandboxPaths() ([]string, []string) {

	if set == nil {
		return []string{}, []string{}
	} else {

		read_paths := make([]string, 0)
		write_paths := make([]string, 0)

		set.mutex.RLock()

		for _, capability := range set.capabilities {

			if capability.Type == CapabilityFilesRead {
				read_paths = append(read_paths, capability.Scope)
			} else if capability.Type == CapabilityFilesWrite {
				read_paths = append(read_paths, capability.Scope)
				write_paths = append(write_paths, capability.Scope)
			}

		}

		set.mutex.RUnlock()

		return UniquePaths(read_paths), UniquePaths(write_paths)

	}
}

func (set *Set) Missing(capabilities []Capability) ([]Capability, error) {

	missing := make([]Capability, 0, len(capabilities))
	seen := make(map[string]bool, len(capabilities))

	for _, capability := range capabilities {

		normalized, err0 := Normalize(capability)

		if err0 == nil {

			if set.Covers(normalized) == false && seen[normalized.Key()] == false {
				seen[normalized.Key()] = true
				missing = append(missing, normalized)
			}

		} else {
			return nil, err0
		}
	}

	sort.Slice(missing, func(a int, b int) bool {
		return missing[a].Key() < missing[b].Key()
	})

	return missing, nil

}

func (set *Set) Clear() {

	if set == nil {
		return
	}

	set.mutex.Lock()
	set.capabilities = make([]Capability, 0)
	set.mutex.Unlock()

}
