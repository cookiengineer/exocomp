package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func WriteFinding(playground string, finding types.Finding) error {

	if finding.Validate() == true {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "findings", finding.ID+".json"))

		if err0 == nil {

			buffer, err1 := json.MarshalIndent(finding, "", "\t")

			if err1 == nil {

				err2 := os.WriteFile(resolved, buffer, 0666)

				if err2 == nil {
					return nil
				} else {
					return fmt.Errorf("findings.WriteFinding: Cannot persist finding \"%s\".", finding.ID)
				}

			} else {
				return fmt.Errorf("findings.WriteFinding: Cannot serialize finding \"%s\".", finding.ID)
			}

		} else {
			return fmt.Errorf("findings.WriteFinding: Cannot resolve finding \"%s\".", finding.ID)
		}

	} else {
		return fmt.Errorf("findings.WriteFinding: Finding \"%s\" is invalid.", finding.ID)
	}

}
