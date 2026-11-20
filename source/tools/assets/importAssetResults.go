package assets

import "exocomp/types"

func importAssetResults(tool *Assets, parsed []types.Asset, invalid int) (int, int, int, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		created := 0
		updated := 0

		for _, asset := range parsed {

			_, is_new, err1 := upsertAsset(tool, asset)

			if err1 == nil {

				if is_new == true {
					created++
				} else {
					updated++
				}

			} else {
				invalid++
			}

		}

		if err2 := writeAssetsIndex(tool); err2 == nil {
			return created, updated, invalid, nil
		} else {
			return 0, 0, 0, err2
		}

	} else {
		return 0, 0, 0, err
	}

}
