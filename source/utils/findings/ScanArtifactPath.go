package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "path/filepath"

func ScanArtifactPath(playground string, scan_id string) (string, error) {

	if types.IsScanID(scan_id) == true {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "scans", scan_id, "results.jsonl"))

		if err0 == nil {
			return resolved, nil
		} else {
			return "", fmt.Errorf("Cannot resolve scan artifact.")
		}

	} else {
		return "", fmt.Errorf("Invalid scan ID \"%s\".", scan_id)
	}
}
