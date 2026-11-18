package assets

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"
import "sort"

func WriteIndex(playground string, entries []types.Asset) error {

	if playground == "" {
		return fmt.Errorf("Invalid playground.")
	}

	resolved, err0 := utils_sandbox.ResolvePath(playground, filepath.Join(".exocomp", "assets", "index.json"))

	if err0 != nil {
		return fmt.Errorf("Cannot resolve assets index.")
	}

	sorted := append([]types.Asset{}, entries...)

	sort.Slice(sorted, func(a int, b int) bool {

		if sorted[a].Type == sorted[b].Type {
			return sorted[a].Value < sorted[b].Value
		}

		return sorted[a].Type < sorted[b].Type

	})

	buffer, err1 := json.MarshalIndent(sorted, "", "\t")

	if err1 != nil {
		return fmt.Errorf("Cannot serialize assets index.")
	}

	if err2 := os.WriteFile(resolved, buffer, 0666); err2 != nil {
		return fmt.Errorf("Cannot persist assets index.")
	}

	return nil

}
