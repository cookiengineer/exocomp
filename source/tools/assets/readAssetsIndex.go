package assets

import utils_findings "exocomp/utils/findings"

func readAssetsIndex(tool *Assets) error {

	entries, err := utils_findings.ReadAssetsIndex(tool.Playground)

	if err == nil {

		for id := range tool.contents {
			delete(tool.contents, id)
		}

		for _, asset := range entries {
			tool.contents[asset.ID] = asset
		}

		return nil

	} else {
		return err
	}

}
