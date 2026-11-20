package lint

import "strings"

func hasPlatformSuffix(base string) bool {

	parts := strings.Split(base, "_")

	if len(parts) < 2 {
		return false
	} else {

		known := map[string]bool{
			"aix": true, "android": true, "darwin": true, "dragonfly": true,
			"freebsd": true, "illumos": true, "ios": true, "js": true,
			"linux": true, "netbsd": true, "openbsd": true, "other": true,
			"plan9": true, "solaris": true, "unix": true, "wasip1": true,
			"windows": true,
			"386":     true, "amd64": true, "arm": true, "arm64": true,
			"loong64": true, "mips": true, "mips64": true, "mips64le": true,
			"mipsle": true, "ppc64": true, "ppc64le": true, "riscv64": true,
			"s390x": true, "wasm": true,
		}

		for index, part := range parts {

			if index == 0 {
				continue
			}

			if known[part] == false {
				return false
			}

		}

		return true

	}

}
