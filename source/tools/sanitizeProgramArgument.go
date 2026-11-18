package tools

import utils_capabilities "exocomp/utils/capabilities"
import utils_sandbox "exocomp/utils/sandbox"
import "net"
import "net/url"
import "path/filepath"
import "strings"
import "unicode"

// sanitizeProgramArgument applies sandbox path rules and capability grants to a
// program argument. URLs and CIDR lists are network values, not paths, and are
// passed through unchanged.
func sanitizeProgramArgument(sandbox string, capabilities *utils_capabilities.Set, argument string) (string, error) {

	if isNetworkArgument(argument) == true {
		return argument, nil
	}

	if strings.ContainsAny(argument, "/\\") == true {

		resolved, err0 := utils_sandbox.SanitizePath(sandbox, argument)

		if err0 == nil {

			return resolved, nil

		} else if capabilities != nil && filepath.IsAbs(argument) == true {

			canonical, err1 := filepath.EvalSymlinks(argument)

			if err1 != nil {

				parent, err2 := filepath.EvalSymlinks(filepath.Dir(argument))

				if err2 == nil {
					canonical = filepath.Join(parent, filepath.Base(argument))
				}

			}

			if canonical != "" {

				if capabilities.CanRead(canonical) == true || capabilities.CanWrite(canonical) == true {
					return canonical, nil
				}

			}

		}

		return "", err0

	} else {

		_, err3 := utils_sandbox.CheckPath(sandbox, argument)

		if err3 == nil {
			return argument, nil
		} else {
			return "", err3
		}

	}

}

// isNetworkArgument reports whether argument is an absolute URL or a
// comma-separated CIDR list. Scanner arguments such as https://target/ or
// 10.0.0.0/8 must not be rewritten as sandbox-relative paths.
func isNetworkArgument(argument string) bool {

	if argument == "" || strings.TrimSpace(argument) != argument {
		return false
	}

	for _, character := range argument {

		if unicode.IsSpace(character) == true || unicode.IsControl(character) == true {
			return false
		}

	}

	if parsed, err := url.Parse(argument); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return true
	}

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
