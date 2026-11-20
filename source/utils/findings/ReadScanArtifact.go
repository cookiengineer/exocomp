package findings

import "fmt"
import "os"

func ReadScanArtifact(playground string, scan_id string) ([]byte, error) {

	resolved, err0 := ScanArtifactPath(playground, scan_id)

	if err0 == nil {

		data, err1 := os.ReadFile(resolved)

		if err1 != nil {

			if os.IsNotExist(err1) == false {
				return nil, fmt.Errorf("Cannot read scan artifact.")
			} else {
				return nil, fmt.Errorf("Scan artifact \"%s\" does not exist.", scan_id)
			}
		}

		return data, nil

	} else {
		return nil, err0
	}
}
