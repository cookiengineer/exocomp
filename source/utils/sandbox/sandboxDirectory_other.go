//go:build !linux && !darwin

package sandbox

func sandboxDirectory() string {
	return "sandbox"
}
