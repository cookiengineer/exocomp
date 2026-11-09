//go:build linux

package sandbox

import utils_landlock "exocomp/utils/landlock"

// Execute applies the platform sandbox and replaces the process with the target.
func Execute(arguments []string) int {
	return utils_landlock.Execute(arguments)
}
