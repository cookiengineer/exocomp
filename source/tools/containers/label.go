package containers

import "crypto/sha256"
import "fmt"
import "path/filepath"

func (tool *Containers) label() string {

	sandbox, err0 := filepath.Abs(tool.Sandbox)

	if err0 == nil {

		resolved, err1 := filepath.EvalSymlinks(sandbox)

		if err1 == nil {
			sandbox = resolved
		}

	} else {
		sandbox = filepath.Clean(tool.Sandbox)
	}

	sum := sha256.Sum256([]byte(sandbox))

	return fmt.Sprintf("%x", sum[0:8])

}
