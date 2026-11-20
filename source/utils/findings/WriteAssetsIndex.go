package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"
import "sort"

// WriteAssetsIndex persists the asset inventory in the playground.
func WriteAssetsIndex(playground string, entries []types.Asset) error {

	if playground != "" {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "assets", "index.json"))

		if err0 == nil {

			sorted := append([]types.Asset{}, entries...)

			sort.Slice(sorted, func(a int, b int) bool {

				if sorted[a].Type == sorted[b].Type {
					return sorted[a].Value < sorted[b].Value
				} else {
					return sorted[a].Type < sorted[b].Type
				}
			})

			buffer, err1 := json.MarshalIndent(sorted, "", "\t")

			if err1 == nil {

				if err2 := os.WriteFile(resolved, buffer, 0666); err2 != nil {
					return fmt.Errorf("Cannot persist assets index.")
				}

				return nil

			} else {
				return fmt.Errorf("Cannot serialize assets index.")
			}

		} else {
			return fmt.Errorf("Cannot resolve assets index.")
		}

	} else {
		return fmt.Errorf("Invalid playground.")
	}
}
