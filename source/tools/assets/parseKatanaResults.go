package assets

import "exocomp/types"
import "encoding/json"
import "strconv"

func parseKatanaResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err == nil {

		result := make([]types.Asset, 0, len(lines))

		for _, line := range lines {

			parsed := katanaImportLine{}

			if json.Unmarshal(line, &parsed) != nil {
				invalid++
				continue
			}

			if parsed.Request.Endpoint == "" {
				invalid++
				continue
			}

			asset := newImportedAsset(types.AssetTypeEndpoint, parsed.Request.Endpoint, "katana", in_scope)
			asset.URL = parsed.Request.Endpoint
			asset.Endpoint = parsed.Request.Endpoint
			asset.Domain = hostFromValue(parsed.Request.Endpoint)
			asset.Props = map[string]string{}

			if parsed.Request.Method != "" {
				asset.Props["method"] = parsed.Request.Method
			}

			if parsed.Response.StatusCode > 0 {
				asset.Props["status_code"] = strconv.Itoa(parsed.Response.StatusCode)
			}

			if parsed.Request.Tag != "" {
				asset.Props["tag"] = parsed.Request.Tag
			}

			result = append(result, asset)

		}

		return result, invalid, nil

	} else {
		return nil, 0, err
	}
}
