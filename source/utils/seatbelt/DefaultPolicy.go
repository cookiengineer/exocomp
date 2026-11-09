package seatbelt

import "path/filepath"

// DefaultPolicy returns the default macOS sandbox access rules for target.
func DefaultPolicy(sandbox string, target string, extra_read []string) Policy {

	read_write := []string{
		sandbox,
	}

	read_only := []string{
		"/usr",
		"/usr/local",
		"/bin",
		"/sbin",
		"/System",
		"/Library",
		"/opt",
		"/etc",
		"/private/etc",
		"/private/var/run",
		"/Applications",
		"/nix/store",
		filepath.Dir(target),
	}

	resolved_target, err0 := filepath.EvalSymlinks(target)

	if err0 == nil {
		read_only = append(read_only, filepath.Dir(resolved_target))
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
