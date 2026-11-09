//go:build !darwin

package seatbelt

// Supported reports that seatbelt is unavailable on this platform.
func Supported() bool {
	return false
}
