package findings

import "fmt"
import "os"

func WriteScanArtifact(playground string, scan_id string, data []byte) error {

	resolved, err0 := ScanArtifactPath(playground, scan_id)

	if err0 == nil {

		if err1 := os.WriteFile(resolved, data, 0644); err1 == nil {
			return nil
		} else {
			return fmt.Errorf("findings.WriteScanArtifact: Cannot persist scan artifact.")
		}
	} else {
		return err0
	}

}
