package sandbox

import "fmt"
import "path/filepath"

func absolutePath(path string) (string, error) {

	absolute, err := filepath.Abs(path)

	if err == nil {
		return absolute, nil
	} else {
		return "", fmt.Errorf("sandbox.AbsolutePath: Cannot resolve helper path.")
	}

}
