package assets

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Assets) Neighbors(id string, arguments utils_chat.Arguments) (string, error) {

	if types.IsAssetID(id) == false {
		return "", fmt.Errorf("assets.Neighbors: Invalid asset ID \"%s\".", id)
	}

	direction := "both"

	if value, ok := arguments.String("direction"); ok == true {
		direction = strings.ToLower(strings.TrimSpace(value))
	}

	if direction != "parent" && direction != "children" && direction != "both" {
		return "", fmt.Errorf("assets.Neighbors: Invalid parameter \"direction\".")
	}

	limit := assetLimit(arguments)

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	asset, err0 := getAsset(tool, id)

	if err0 == nil {

		lines := make([]string, 0)

		if (direction == "parent" || direction == "both") && asset.ParentID != "" {

			if parent, ok := tool.contents[asset.ParentID]; ok == true {
				lines = append(lines, "Parent: "+formatAssetSummary(parent))
			}

		}

		if direction == "children" || direction == "both" {

			children := make([]types.Asset, 0)

			for _, candidate := range tool.contents {

				if candidate.ParentID == id {
					children = append(children, candidate)
				}

			}

			sortAssets(children)

			if len(children) > limit {
				children = children[0:limit]
			}

			if len(children) > 0 {
				lines = append(lines, fmt.Sprintf("Children: %d", len(children)))
			}

			for _, child := range children {
				lines = append(lines, formatAssetSummary(child))
			}

		}

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("assets.Neighbors: %s", id))

		if len(lines) == 0 {
			result = append(result, "No neighbors.")
		} else {
			result = append(result, lines...)
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("assets.Neighbors: %s", err0.Error())
	}

}
