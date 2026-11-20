package assets

import "exocomp/types"
import "encoding/json"
import "fmt"
import "strings"

func (tool *Assets) Get(id string) (string, error) {

	if types.IsAssetID(id) == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		asset, err0 := getAsset(tool, id)

		if err0 == nil {

			buffer, err1 := json.MarshalIndent(asset, "", "\t")

			if err1 == nil {
				return strings.Join([]string{fmt.Sprintf("assets.Get: %s", id), string(buffer)}, "\n"), nil
			} else {
				return "", fmt.Errorf("assets.Get: Cannot serialize asset \"%s\".", id)
			}

		} else {
			return "", fmt.Errorf("assets.Get: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Get: Invalid asset ID \"%s\".", id)
	}

}
