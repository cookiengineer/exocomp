package assets

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"

func (tool *Assets) Update(id string, arguments utils_chat.Arguments) (string, error) {

	if types.IsAssetID(id) == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		asset, err0 := getAsset(tool, id)

		if err0 == nil {
			return updateAsset(tool, id, asset, arguments)
		} else {
			return "", fmt.Errorf("assets.Update: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Update: Invalid asset ID \"%s\".", id)
	}

}
