package findings

import "fmt"
import "os"

func ReadScanArtifact(playground string, scan_id string) ([]byte, error) {

	resolved, err0 := ScanArtifactPath(playground, scan_id)

	if err0 != nil {
		return nil, err0
	}

	data, err1 := os.ReadFile(resolved)

	if err1 != nil {

		if os.IsNotExist(err1) == true {
			return nil, fmt.Errorf("Scan artifact \"%s\" does not exist.", scan_id)
		}

		return nil, fmt.Errorf("Cannot read scan artifact.")

	}

	return data, nil

}
