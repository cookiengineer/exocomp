package landlock

import "path/filepath"

// DefaultPolicy returns the default sandbox access rules for target.
func DefaultPolicy(sandbox string, target string, extra_read []string) Policy {

	read_write := []string{
		sandbox,
	}

	read_only := []string{
		"/usr",
		"/usr/local",
		"/bin",
		"/sbin",
		"/lib",
		"/lib32",
		"/lib64",
		"/libx32",
		"/etc",
		"/opt",
		"/run",
		"/nix/store",
		"/proc",
		"/sys",
		filepath.Dir(target),
	}

	read_only = append(read_only, extra_read...)

	devices := []string{
		"/dev/null",
		"/dev/zero",
		"/dev/full",
		"/dev/random",
		"/dev/urandom",
		"/dev/tty",
	}

	return Policy{
		ReadWrite: uniquePaths(read_write),
		ReadOnly:  uniquePaths(read_only),
		Devices:   uniquePaths(devices),
	}

}
