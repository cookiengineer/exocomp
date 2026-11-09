//go:build darwin

package sandbox

import utils_seatbelt "exocomp/utils/seatbelt"

// Execute applies the platform sandbox and replaces the process with the target.
func Execute(arguments []string) int {
	return utils_seatbelt.Execute(arguments)
}
