package findings

import "exocomp/types"
import "strings"

// NormalizeAssetValue returns a stable, case-insensitive value for asset
// deduplication.
func NormalizeAssetValue(asset_type types.AssetType, value string) string {

	value = strings.TrimSpace(value)

	switch asset_type {
	case types.AssetTypeIP, types.AssetTypeCIDR, types.AssetTypeDomain, types.AssetTypeSubdomain:
		return strings.ToLower(value)
	case types.AssetTypeURL:
		return normalizeAssetURL(value)
	default:
		return value
	}

}
