package assets

import "exocomp/types"
import "strings"

func collectExportAssets(tool *Assets, type_filter string, has_type bool, query_filter string, has_query bool, in_scope_filter bool, has_in_scope bool, tested_filter bool, has_tested bool, limit int) ([]types.Asset, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		matches := make([]types.Asset, 0)

		for _, asset := range tool.contents {

			if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
				continue
			} else if has_query == true && assetMatchesQuery(asset, query_filter) == false {
				continue
			} else if has_in_scope == true && asset.InScope != in_scope_filter {
				continue
			} else if has_tested == true && asset.Tested != tested_filter {
				continue
			}

			matches = append(matches, asset)

		}

		sortAssets(matches)

		if len(matches) > limit {
			matches = matches[0:limit]
		}

		return matches, nil

	} else {
		return nil, err
	}

}
