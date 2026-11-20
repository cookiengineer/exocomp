package shell

import "net"
import "net/url"
import "strings"
import "unicode"

// IsNetworkArgument reports whether argument is an absolute URL or a
// comma-separated CIDR list. Scanner arguments such as https://target/ or
// 10.0.0.0/8 must not be rewritten as sandbox-relative paths.
func IsNetworkArgument(argument string) bool {

	if argument != "" && strings.TrimSpace(argument) == argument {

		for _, character := range argument {

			if unicode.IsSpace(character) == true || unicode.IsControl(character) == true {
				return false
			}

		}

		if parsed, err := url.Parse(argument); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			return true
		} else {

			parts := strings.Split(argument, ",")
			valid := len(parts) > 0

			for _, part := range parts {

				if _, _, err := net.ParseCIDR(strings.TrimSpace(part)); err != nil {

					valid = false
					break

				}

			}

			return valid

		}

	} else {
		return false
	}
}
