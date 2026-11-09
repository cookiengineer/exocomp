//go:build linux

package landlock

// Supported reports whether the running kernel supports landlock.
func Supported() bool {

	abi, err := landlockABI()

	if err == nil && abi >= 1 {
		return true
	} else {
		return false
	}

}
