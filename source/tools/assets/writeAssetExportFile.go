package assets

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"

func writeAssetExportFile(sandbox string, path string, data []byte) (string, string, error) {

	resolved, err0 := utils_sandbox.ResolvePath(sandbox, path)

	if err0 == nil {

		if err1 := os.WriteFile(resolved, data, 0644); err1 == nil {

			display, err2 := utils_sandbox.SanitizePath(sandbox, resolved)

			if err2 == nil {
				return resolved, display, nil
			} else {
				return "", "", fmt.Errorf("Cannot resolve export path.")
			}

		} else {
			return "", "", fmt.Errorf("Cannot write export file.")
		}

	} else {
		return "", "", fmt.Errorf("Cannot resolve export path.")
	}

}
