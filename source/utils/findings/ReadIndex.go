package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func ReadIndex(playground string) ([]types.FindingSummary, error) {

	if playground != "" {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "findings", "index.json"))

		if err0 == nil {

			buffer, err1 := os.ReadFile(resolved)

			if err1 == nil {

				summaries := make([]types.FindingSummary, 0)
				err2 := json.Unmarshal(buffer, &summaries)

				if err2 == nil {
					return summaries, nil
				} else {
					return nil, fmt.Errorf("findings.ReadIndex: Cannot parse findings index.")
				}

			} else if os.IsNotExist(err1) == true {
				return make([]types.FindingSummary, 0), nil
			} else {
				return nil, fmt.Errorf("findings.ReadIndex: Cannot read findings index.")
			}

		} else {
			return nil, fmt.Errorf("findings.ReadIndex: Cannot resolve findings index.")
		}

	} else {
		return nil, fmt.Errorf("findings.ReadIndex: Invalid playground.")
	}

}
