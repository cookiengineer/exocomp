package assets

import "exocomp/types"

func storeAsset(tool *Assets, asset types.Asset) (types.Asset, bool, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		stored, created, err1 := upsertAsset(tool, asset)

		if err1 == nil {
			return stored, created, writeAssetsIndex(tool)
		} else {
			return types.Asset{}, false, err1
		}

	} else {
		return types.Asset{}, false, err
	}

}
