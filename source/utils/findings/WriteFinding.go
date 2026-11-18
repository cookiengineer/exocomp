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
					return fmt.Errorf("Cannot persist finding \"%s\".", finding.ID)
				}

			} else {
				return fmt.Errorf("Cannot serialize finding \"%s\".", finding.ID)
			}

		} else {
			return fmt.Errorf("Cannot resolve finding \"%s\".", finding.ID)
		}

	} else {
		return fmt.Errorf("Finding \"%s\" is invalid.", finding.ID)
	}

}
