package testutil

import "path/filepath"
import "runtime"

// sourceRoot returns the absolute path of the source folder.
func sourceRoot() string {

	_, filename, _, ok := runtime.Caller(0)

	if ok == true {
		return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	} else {
		return ""
	}

}
