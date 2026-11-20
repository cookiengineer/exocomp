package assets

import "exocomp/types"
import "fmt"

func getAsset(tool *Assets, id string) (types.Asset, error) {

	if err := readAssetsIndex(tool); err == nil {

		asset, ok := tool.contents[id]

		if ok == true {
			return asset, nil
		} else {
			return types.Asset{}, fmt.Errorf("Asset \"%s\" does not exist.", id)
		}

	} else {
		return types.Asset{}, err
	}

}
