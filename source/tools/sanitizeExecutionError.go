package tools

import "errors"
import "fmt"
import "io/fs"
import "strings"

func sanitizeExecutionError(namespace string, method string, kind string, name string, result string, err error) (string, error) {

	title := strings.ToUpper(kind[0:1]) + kind[1:]

	if errors.Is(err, fs.ErrPermission) {
		return "", fmt.Errorf("%s.%s: Invalid %s \"%s\": Permission denied.", namespace, method, kind, name)
	} else if errors.Is(err, fs.ErrNotExist) || strings.Contains(err.Error(), "executable file not found") {
		return "", fmt.Errorf("%s.%s: Invalid %s \"%s\": Program doesn't exist.", namespace, method, kind, name)
	} else {
		return result, fmt.Errorf("%s.%s: %s \"%s\" exited with an error.", namespace, method, title, name)
	}

}
