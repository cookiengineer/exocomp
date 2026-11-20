package landlock

import utils_slices "exocomp/utils/slices"
import "path/filepath"

// DefaultPolicy returns the default sandbox access rules for target.
func DefaultPolicy(sandbox string, target string, extra_read []string, extra_write []string) Policy {

	read_write := []string{
		sandbox,
	}

	read_write = append(read_write, extra_write...)

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
		"/proc/self",
		"/proc/thread-self",
		"/proc/cpuinfo",
		"/proc/filesystems",
		"/proc/loadavg",
		"/proc/meminfo",
		"/proc/stat",
		"/proc/sys",
		"/proc/uptime",
		"/proc/version",
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
		ReadWrite: utils_slices.UniqueStrings(read_write),
		ReadOnly:  utils_slices.UniqueStrings(read_only),
		Devices:   utils_slices.UniqueStrings(devices),
	}

}
