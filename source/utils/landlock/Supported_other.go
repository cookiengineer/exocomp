//go:build !linux

package landlock

// Supported reports that landlock is unavailable on this platform.
func Supported() bool {
	return false
}
