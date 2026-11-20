package seatbelt

import "os"
import "path/filepath"

func resolvePaths(paths []string) []resolvedPath {

	result := make([]resolvedPath, 0, len(paths))
	seen := make(map[string]bool, len(paths))

	for _, path := range paths {

		if path == "" {
			continue
		} else {

			cleaned := filepath.Clean(path)

			if seen[cleaned] == true {
				continue
			} else {

				seen[cleaned] = true

				info, err0 := os.Stat(cleaned)

				if err0 == nil {

					resolved, err1 := filepath.EvalSymlinks(cleaned)

					if err1 != nil {
						resolved = cleaned
					}

					result = append(result, resolvedPath{
						original: cleaned,
						resolved: resolved,
						isDir:    info.IsDir(),
						exists:   true,
					})

				} else {

					resolved, err2 := filepath.EvalSymlinks(filepath.Dir(cleaned))

					if err2 == nil {
						resolved = filepath.Join(resolved, filepath.Base(cleaned))
					} else {
						resolved = cleaned
					}

					result = append(result, resolvedPath{
						original: cleaned,
						resolved: resolved,
						isDir:    true,
						exists:   false,
					})

				}

			}

		}
	}

	return result

}
