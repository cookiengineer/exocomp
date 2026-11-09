package seatbelt

import "fmt"
import "os"
import "path/filepath"

type resolvedPath struct {
	original string
	resolved string
	isDir    bool
	exists   bool
}

func resolvePaths(paths []string) []resolvedPath {

	result := make([]resolvedPath, 0, len(paths))
	seen := make(map[string]bool, len(paths))

	for _, path := range paths {

		if path == "" {
			continue
		}

		cleaned := filepath.Clean(path)

		if seen[cleaned] == true {
			continue
		}

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

	return result

}

func pathFilter(resolved_path resolvedPath) string {

	if resolved_path.exists == true {

		matcher := "subpath"

		if resolved_path.isDir == false {
			matcher = "literal"
		}

		if resolved_path.original == resolved_path.resolved {
			return fmt.Sprintf("(%s %s)", matcher, quoteSBPL(resolved_path.original))
		} else {
			return fmt.Sprintf(
				"(require-any (%s %s) (%s %s))",
				matcher, quoteSBPL(resolved_path.original),
				matcher, quoteSBPL(resolved_path.resolved),
			)
		}

	} else {

		if resolved_path.original == resolved_path.resolved {
			return fmt.Sprintf(
				"(require-any (literal %s) (subpath %s))",
				quoteSBPL(resolved_path.original),
				quoteSBPL(resolved_path.original),
			)
		} else {
			return fmt.Sprintf(
				"(require-any (literal %s) (subpath %s) (literal %s) (subpath %s))",
				quoteSBPL(resolved_path.original),
				quoteSBPL(resolved_path.original),
				quoteSBPL(resolved_path.resolved),
				quoteSBPL(resolved_path.resolved),
			)
		}

	}

}

func pathFilters(paths []resolvedPath) string {

	result := ""

	for index, resolved_path := range paths {

		if index > 0 {
			result = result + " "
		}

		result = result + pathFilter(resolved_path)

	}

	return result

}
