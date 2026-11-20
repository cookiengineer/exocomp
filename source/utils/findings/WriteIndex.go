package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func WriteIndex(playground string, summaries []types.FindingSummary) error {

	if playground != "" {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "findings", "index.json"))

		if err0 == nil {

			sorted := append([]types.FindingSummary{}, summaries...)
			SortSummaries(sorted)

			buffer, err1 := json.MarshalIndent(sorted, "", "\t")

			if err1 == nil {

				err2 := os.WriteFile(resolved, buffer, 0666)

				if err2 == nil {
					return nil
				} else {
					return fmt.Errorf("findings.WriteIndex: Cannot persist findings index.")
				}

			} else {
				return fmt.Errorf("findings.WriteIndex: Cannot serialize findings index.")
			}

		} else {
			return fmt.Errorf("findings.WriteIndex: Cannot resolve findings index.")
		}

	} else {
		return fmt.Errorf("findings.WriteIndex: Invalid playground.")
	}

}
