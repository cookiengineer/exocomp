package assets

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Assets) Search(arguments utils_chat.Arguments) (string, error) {

	query, ok := arguments.String("query")

	if ok == false || strings.TrimSpace(query) == "" {
		return "", fmt.Errorf("assets.Search: Invalid parameter \"query\" is not a non-empty string.")
	}

	type_filter, has_type := arguments.String("type")
	limit := assetLimit(arguments)

	if has_type == true && types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == false {
		return "", fmt.Errorf("assets.Search: Invalid parameter \"type\".")
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		matches := make([]types.Asset, 0)

		for _, asset := range tool.contents {

			if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
				continue
			} else if assetMatchesQuery(asset, query) == false {
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
		result = append(result, fmt.Sprintf("assets.Search: %d matches (showing %d).", total, len(lines)))
		result = append(result, lines...)

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("assets.Search: %s", err.Error())
	}

}
