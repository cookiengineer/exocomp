package assets

import "exocomp/types"
import "fmt"
import "sort"
import "strings"

func (tool *Assets) Coverage() (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		counts := make(map[types.AssetType][3]int)

		for _, asset := range tool.contents {

			count := counts[asset.Type]

			if asset.InScope == true {
				count[0]++
			}

			if asset.Tested == true {
				count[1]++
			}

			count[2]++
			counts[asset.Type] = count

		}

		types_list := make([]string, 0, len(counts))

		for asset_type := range counts {
			types_list = append(types_list, string(asset_type))
		}

		sort.Strings(types_list)

		lines := make([]string, 0, len(types_list))

		for _, name := range types_list {

			count := counts[types.AssetType(name)]
			lines = append(lines, fmt.Sprintf("%s: %d assets, %d in scope, %d tested", name, count[2], count[0], count[1]))

		}

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("assets.Coverage: %d assets across %d types.", len(tool.contents), len(types_list)))
		result = append(result, lines...)

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("assets.Coverage: %s", err.Error())
	}

}
