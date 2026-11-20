package assets

import "exocomp/types"
import "encoding/json"

func parseUrlfinderResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err == nil {

		result := make([]types.Asset, 0, len(lines))

		for _, line := range lines {

			parsed := urlfinderImportLine{}

			if json.Unmarshal(line, &parsed) != nil {
				invalid++
				continue
			}

			if parsed.URL == "" {
				invalid++
				continue
			}

			asset := newImportedAsset(types.AssetTypeURL, parsed.URL, "urlfinder", in_scope)
			asset.URL = parsed.URL
			asset.Endpoint = parsed.URL
			asset.Domain = hostFromValue(parsed.URL)
			asset.Props = map[string]string{}

			if parsed.Source != "" {
				asset.Props["source"] = parsed.Source
			}

			result = append(result, asset)

		}

		return result, invalid, nil

	} else {
		return nil, 0, err
	}
}
