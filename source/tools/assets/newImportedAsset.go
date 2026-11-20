package assets

import "exocomp/types"
import utils_findings "exocomp/utils/findings"

func newImportedAsset(asset_type types.AssetType, value string, source string, in_scope bool) types.Asset {

	return types.Asset{
		Type:       asset_type,
		Value:      utils_findings.LimitText(value, assetsMaxValue),
		Source:     source,
		Confidence: types.FindingConfidenceObserved,
		InScope:    in_scope,
	}
}
