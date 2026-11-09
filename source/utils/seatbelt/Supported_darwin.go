//go:build darwin

package seatbelt

// Supported reports whether the macOS seatbelt sandbox can be applied.
func Supported() bool {

	_, err := sandboxExecutable()

	if err == nil {
		return true
	} else {
		return false
	}

}
