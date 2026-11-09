//go:build !linux && !darwin

package sandbox

// Supported reports that sandboxing is unavailable on this platform.
func Supported() bool {
	return false
}
