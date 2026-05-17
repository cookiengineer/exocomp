//go:build linux || darwin || freebsd || openbsd || netbsd

package fs

import "golang.org/x/sys/unix"
import "errors"
import "path/filepath"

func AvailableSpace(path string) uint64 {

	result := uint64(0)

	stat := unix.Statfs_t{}
	err0 := unix.Statfs(path, &stat)

	if err0 == nil {

		result = uint64(stat.Bavail * uint64(stat.Bsize))

	} else {

		for {

			parent := filepath.Dir(path)

			if parent == path {
				break
			}

			stat := unix.Statfs_t{}
			err1 := unix.Statfs(parent, &stat)

			if err1 == nil {

				result = uint64(stat.Bavail) * uint64(stat.Bsize)
				break

			} else if errors.Is(err1, unix.ENOENT) {

				path = parent
				continue

			} else {
				break
			}

		}

	}

	return result

}

