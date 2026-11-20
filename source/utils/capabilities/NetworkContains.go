package capabilities

import "net"

// NetworkContains reports whether an IP matches an exact IP or CIDR scope.
func NetworkContains(scope string, ip net.IP) bool {

	if parsed := net.ParseIP(scope); parsed != nil {
		return parsed.Equal(ip)
	} else if _, network, err := net.ParseCIDR(scope); err == nil {
		return network.Contains(ip)
	} else {
		return false
	}
}
