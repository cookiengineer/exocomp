//go:build linux

package sandbox

import utils_landlock "exocomp/utils/landlock"

// Supported reports whether sandboxing is available on this system.
func Supported() bool {
	return utils_landlock.Supported()
}
