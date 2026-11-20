package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "path/filepath"

func exportScanData(tool *Findings, arguments utils_chat.Arguments) ([]byte, string, string, error) {

	id, has_id := arguments.String("id")

	if has_id == true && types.IsScanID(id) == true {

		buffer, err0 := utils_findings.ReadScanArtifact(tool.Playground, id)

		if err0 == nil {
			return buffer, filepath.Join("work", "scans", id+"-results.jsonl"), "application/x-ndjson", nil
		} else {
			return nil, "", "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
	}

}
