package assets

import "exocomp/types"
import "fmt"
import "strings"

func formatAssetSummary(asset types.Asset) string {

	flags := make([]string, 0)

	if asset.InScope == true {
		flags = append(flags, "in_scope")
	}

	if asset.Tested == true {
		flags = append(flags, "tested")
	}

	suffix := ""

	if len(flags) > 0 {
		suffix = " (" + strings.Join(flags, ", ") + ")"
	}

	return fmt.Sprintf("%s [%s] %s%s", asset.ID, asset.Type, asset.Value, suffix)

}
