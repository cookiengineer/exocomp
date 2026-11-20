package assets

import "exocomp/types"
import utils_findings "exocomp/utils/findings"

func writeAssetsIndex(tool *Assets) error {

	entries := make([]types.Asset, 0, len(tool.contents))

	for _, asset := range tool.contents {
		entries = append(entries, asset)
	}

	return utils_findings.WriteAssetsIndex(tool.Playground, entries)

}
