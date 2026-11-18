package findings

import "fmt"
import "os"

func WriteScanArtifact(playground string, scan_id string, data []byte) error {

	resolved, err0 := ScanArtifactPath(playground, scan_id)

	if err0 != nil {
		return err0
	}

	if err1 := os.WriteFile(resolved, data, 0644); err1 != nil {
		return fmt.Errorf("Cannot persist scan artifact.")
	}

	return nil

}
