package assets

import "exocomp/types"
import "net/url"
import "strings"

// NormalizeValue returns a stable, case-insensitive value for deduplication.
func NormalizeValue(asset_type types.AssetType, value string) string {

	value = strings.TrimSpace(value)

	switch asset_type {
	case types.AssetTypeIP, types.AssetTypeCIDR, types.AssetTypeDomain, types.AssetTypeSubdomain:
		return strings.ToLower(value)
	case types.AssetTypeURL:
		return normalizeURL(value)
	default:
		return value
	}

}

func normalizeURL(value string) string {

	parsed, err := url.Parse(value)

	if err == nil && parsed.Scheme != "" && parsed.Host != "" {

		parsed.Scheme = strings.ToLower(parsed.Scheme)
		parsed.Host = strings.ToLower(parsed.Host)

		return parsed.String()

	}

	return value

}
