package fs

import "os"
import "path/filepath"

func HasExocompSession(folder string) bool {

	tmp         := filepath.Join(folder, ".exocomp", "session.json")
	stat, err := os.Stat(tmp)

	if err == nil && stat.IsDir() {
		return true
	}

	return false

}

