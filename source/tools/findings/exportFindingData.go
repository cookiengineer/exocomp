package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "encoding/json"
import "fmt"
import "path/filepath"

func exportFindingData(tool *Findings, arguments utils_chat.Arguments) ([]byte, string, string, error) {

	id, has_id := arguments.String("id")

	if has_id == true && types.IsFindingID(id) == true {

		finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

		if err0 == nil {

			buffer, err1 := json.MarshalIndent(finding, "", "\t")

			if err1 == nil {
				return buffer, filepath.Join("work", "findings", id+".json"), "application/json", nil
			} else {
				return nil, "", "", fmt.Errorf("findings.Export: Cannot serialize finding \"%s\".", id)
			}

		} else {
			return nil, "", "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
	}

}
