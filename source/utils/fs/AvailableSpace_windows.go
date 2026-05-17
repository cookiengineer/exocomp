//go:build windows

package fs

import "golang.org/x/sys/windows"

func AvailableSpace(path string) uint64 {

	result := uint64(0)

	path_pointer, err0 := windows.UTF16PtrFromString(path)

	if err0 == nil {

		free_bytes_available := uint64(0)

		err1 := windows.GetDiskFreeSpaceEx(
			path_pointer,
			&free_bytes_available,
			nil,
			nil,
		)

		if err1 == nil {

			result = uint64(free_bytes_available)

		} else {

			for {

				parent := filepath.Dir(path)

				if parent == path {
					break
				}

				parent_pointer, err2 := windows.UTF16PtrFromString(parent)

				if err2 == nil {

					err3 := windows.GetDiskFreeSpaceEx(
						parent_pointer,
						&free_bytes_available,
						nil,
						nil,
					)

					if err3 == nil {

						result = uint64(free_bytes_available)
						break

					} else if errors.Is(err3, windows.ERROR_PATH_NOT_FOUND) {

						path = parent
						continue

					} else {
						break
					}

				} else {
					break
				}

			}

		}

	}

	return result

}
