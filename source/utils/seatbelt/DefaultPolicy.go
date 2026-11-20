package seatbelt

import utils_slices "exocomp/utils/slices"
import "path/filepath"

// DefaultPolicy returns the default macOS sandbox access rules for target.
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

	deny := []string{
		filepath.Join(sandbox, ".exocomp"),
	}

	return Policy{
		ReadWrite: utils_slices.UniqueStrings(read_write),
		ReadOnly:  utils_slices.UniqueStrings(read_only),
		Devices:   utils_slices.UniqueStrings(devices),
		Deny:      utils_slices.UniqueStrings(deny),
	}

}
