package assets

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Assets) List(arguments utils_chat.Arguments) (string, error) {

	type_filter, has_type := arguments.String("type")
	query, has_query := arguments.String("query")
	in_scope_filter, has_in_scope := arguments.Bool("in_scope")
	tested_filter, has_tested := arguments.Bool("tested")
	limit := assetLimit(arguments)

	if has_type == true && types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == false {
		return "", fmt.Errorf("assets.List: Invalid parameter \"type\".")
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		matches := make([]types.Asset, 0)

		for _, asset := range tool.contents {

			if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
				continue
			} else if has_query == true && assetMatchesQuery(asset, query) == false {
				continue
			} else if has_in_scope == true && asset.InScope != in_scope_filter {
				continue
			} else if has_tested == true && asset.Tested != tested_filter {
				continue
			}

			matches = append(matches, asset)

		}

		sortAssets(matches)

		total := len(matches)

		if len(matches) > limit {
			matches = matches[0:limit]
		}

		lines := make([]string, 0, len(matches))

		for _, asset := range matches {
			lines = append(lines, formatAssetSummary(asset))
		}

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("assets.List: %d assets (showing %d).", total, len(lines)))
		result = append(result, lines...)

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("assets.List: %s", err.Error())
	}

}
