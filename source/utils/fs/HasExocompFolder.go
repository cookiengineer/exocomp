package fs

import "os"
import "path/filepath"

func HasExocompFolder(folder string) bool {

	tmp         := filepath.Join(folder, ".exocomp")
	stat, err := os.Stat(tmp)

	if err == nil && stat.IsDir() {
		return true
	}

	return false

}

