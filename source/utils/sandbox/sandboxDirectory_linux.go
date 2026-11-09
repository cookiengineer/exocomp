//go:build linux

package sandbox

func sandboxDirectory() string {
	return "landlock"
}
