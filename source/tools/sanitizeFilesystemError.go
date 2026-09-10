package tools

import "errors"
import "fmt"
import "io/fs"

func sanitizeFilesystemError(namespace string, method string, kind string, path string, err error) (string, error) {

	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("%s.%s: %s \"%s\" does not exist.", namespace, method, kind, path)
	} else if errors.Is(err, fs.ErrPermission) {
		return "", fmt.Errorf("%s.%s: %s \"%s\": Permission denied.", namespace, method, kind, path)
	} else {
		return "", fmt.Errorf("%s.%s: Cannot access %s \"%s\".", namespace, method, kind, path)
	}

}
