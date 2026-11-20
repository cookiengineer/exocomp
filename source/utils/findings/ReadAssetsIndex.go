package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

// ReadAssetsIndex reads the asset inventory from the playground.
func ReadAssetsIndex(playground string) ([]types.Asset, error) {

	if playground != "" {

		resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "assets", "index.json"))

		if err0 == nil {

			buffer, err1 := os.ReadFile(resolved)

			if err1 == nil {

				result := make([]types.Asset, 0)

				if err2 := json.Unmarshal(buffer, &result); err2 != nil {
					return nil, fmt.Errorf("Cannot parse assets index.")
				}

				return result, nil

			} else if os.IsNotExist(err1) == true {
				return make([]types.Asset, 0), nil
			} else {
				return nil, fmt.Errorf("Cannot read assets index.")
			}

		} else {
			return nil, fmt.Errorf("Cannot resolve assets index.")
		}

	} else {
		return nil, fmt.Errorf("Invalid playground.")
	}
}
