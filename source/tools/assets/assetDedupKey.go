package assets

import "exocomp/types"
import utils_findings "exocomp/utils/findings"

func assetDedupKey(asset types.Asset) string {
	return string(asset.Type) + "\x00" + utils_findings.NormalizeAssetValue(asset.Type, asset.Value)
}
