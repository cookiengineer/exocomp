package fs

import "fmt"
import "os"

func CreateSandbox(prefix string) (string, error) {

	if prefix == "" {
		prefix = "exocomp"
	}

	folder, err := os.MkdirTemp("", fmt.Sprintf("%s-*", prefix))

	if err == nil {
		return folder, nil
	} else {
		return "", err
	}

}
