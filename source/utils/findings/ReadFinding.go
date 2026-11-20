package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func ReadFinding(playground string, id string) (types.Finding, error) {

	if types.IsFindingID(id) == true {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "findings", id+".json"))

		if err0 == nil {

			buffer, err1 := os.ReadFile(resolved)

			if err1 == nil {

				finding := types.Finding{}
				err2 := json.Unmarshal(buffer, &finding)

				if err2 == nil {

					if finding.Validate() == true {
						return finding, nil
					} else {
						return types.Finding{}, fmt.Errorf("findings.ReadFinding: Finding \"%s\" is invalid.", id)
					}

				} else {
					return types.Finding{}, fmt.Errorf("findings.ReadFinding: Cannot parse finding \"%s\".", id)
				}

			} else if os.IsNotExist(err1) == true {
				return types.Finding{}, fmt.Errorf("findings.ReadFinding: Finding \"%s\" does not exist.", id)
			} else {
				return types.Finding{}, fmt.Errorf("findings.ReadFinding: Cannot read finding \"%s\".", id)
			}

		} else {
			return types.Finding{}, fmt.Errorf("findings.ReadFinding: Cannot resolve finding \"%s\".", id)
		}

	} else {
		return types.Finding{}, fmt.Errorf("findings.ReadFinding: Invalid finding ID \"%s\".", id)
	}

}
