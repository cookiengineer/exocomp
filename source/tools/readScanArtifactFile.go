package tools

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"

func readScanArtifactFile(sandbox string, path string, namespace string) ([]byte, error) {

	resolved, err0 := utils_sandbox.ResolvePath(sandbox, path)

	if err0 == nil {

		stat, err1 := os.Stat(resolved)

		if err1 == nil && stat.Mode().IsRegular() == true {

			if stat.Size() <= maxScanArtifactSize {

				data, err2 := os.ReadFile(resolved)

				if err2 == nil {
					return data, nil
				} else {
					return nil, fmt.Errorf("%s: Cannot read artifact \"%s\".", namespace, path)
				}

			} else {
				return nil, fmt.Errorf("%s: Artifact exceeds the 16 MiB limit.", namespace)
			}

		} else {
			return nil, fmt.Errorf("%s: Artifact \"%s\" does not exist.", namespace, path)
		}

	} else {
		return nil, fmt.Errorf("%s: Cannot resolve artifact \"%s\".", namespace, path)
	}

}
