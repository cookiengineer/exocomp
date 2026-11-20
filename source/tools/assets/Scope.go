package assets

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Assets) Scope(arguments utils_chat.Arguments) (string, error) {

	type_filter, has_type := arguments.String("type")
	value_filter, has_value := arguments.String("value")
	query, has_query := arguments.String("query")
	has_filter := has_type == true || has_value == true || has_query == true
	in_scope := true

	if provided, ok := arguments.Bool("in_scope"); ok == true {
		in_scope = provided
	}

	if has_type == false || types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		if err := readAssetsIndex(tool); err == nil {

			if has_filter == false {

				total := len(tool.contents)
				scoped := 0
				tested := 0

				for _, asset := range tool.contents {

					if asset.InScope == true {
						scoped++
					}

					if asset.Tested == true {
						tested++
					}

				}

				return fmt.Sprintf("assets.Scope: %d assets, %d in scope, %d tested.", total, scoped, tested), nil

			} else {

				updated := 0

				for id, asset := range tool.contents {

					if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
						continue
					} else if has_value == true && strings.EqualFold(asset.Value, strings.TrimSpace(value_filter)) == false {
						continue
					} else if has_query == true && assetMatchesQuery(asset, query) == false {
						continue
					}

					asset.InScope = in_scope
					asset.LastSeen = schemas.NewDatetime()
					tool.contents[id] = asset
					updated++

				}

				if err := writeAssetsIndex(tool); err == nil {
					return fmt.Sprintf("assets.Scope: %d assets marked in_scope=%t.", updated, in_scope), nil
				} else {
					return "", fmt.Errorf("assets.Scope: %s", err.Error())
				}

			}

		} else {
			return "", fmt.Errorf("assets.Scope: %s", err.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Scope: Invalid parameter \"type\".")
	}
}
