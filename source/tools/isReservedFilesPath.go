package tools

import "path/filepath"
import "strings"

// isReservedFilesPath reports whether path points at the internal .exocomp
// folder. The folder is reserved for findings, assets, shell transcripts and
// scanner artifacts, which must be accessed through their own tools.
func isReservedFilesPath(sandbox string, path string) bool {

	cleaned := filepath.Clean(path)

	if filepath.IsAbs(cleaned) == true {

		sandbox_root, err := filepath.Abs(sandbox)

		if err == nil {

			relative, err := filepath.Rel(sandbox_root, cleaned)

			if err == nil && relative != ".." && strings.HasPrefix(relative, ".."+string(filepath.Separator)) == false {
				cleaned = relative
			}

		}

	}

	for _, part := range strings.Split(cleaned, string(filepath.Separator)) {

		if strings.EqualFold(part, ".exocomp") == true {
			return true
		}

	}

	return false

}
