//go:build darwin

package sandbox

import utils_seatbelt "exocomp/utils/seatbelt"

// Supported reports whether sandboxing is available on this system.
func Supported() bool {
	return utils_seatbelt.Supported()
}
