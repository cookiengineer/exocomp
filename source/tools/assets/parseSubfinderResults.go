package assets

import "exocomp/types"
import "encoding/json"
import "strings"

func parseSubfinderResults(data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := subfinderImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.Host == "" {
			invalid++
			continue
		}

		asset_type := types.AssetTypeSubdomain

		if target != "" && strings.EqualFold(parsed.Host, target) == true {
			asset_type = types.AssetTypeDomain
		}

		asset := newImportedAsset(asset_type, parsed.Host, "subfinder", in_scope)
		asset.Domain = parsed.Host
		asset.Props = map[string]string{}

		if parsed.Source != "" {
			asset.Props["source"] = parsed.Source
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}
