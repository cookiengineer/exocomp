package assets

import "exocomp/types"
import utils_strings "exocomp/utils/strings"
import "encoding/json"
import "strconv"
import "strings"

func parseHttpxResults(data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := httpxImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		value := utils_strings.FirstNonEmpty(parsed.URL, parsed.Host, parsed.Input)

		if value == "" {
			invalid++
			continue
		}

		asset_type := types.AssetTypeSubdomain

		if strings.Contains(value, "://") == true {
			asset_type = types.AssetTypeURL
		} else if target != "" && strings.EqualFold(value, target) == true {
			asset_type = types.AssetTypeDomain
		}

		asset := newImportedAsset(asset_type, value, "httpx", in_scope)
		asset.URL = parsed.URL
		asset.Endpoint = parsed.URL
		asset.IP = parsed.IP
		asset.Domain = hostFromValue(value)
		asset.Technology = append([]string{}, parsed.Tech...)
		asset.Props = map[string]string{}

		if parsed.StatusCode > 0 {
			asset.Props["status_code"] = strconv.Itoa(parsed.StatusCode)
		}

		if parsed.Title != "" {
			asset.Props["title"] = parsed.Title
		}

		if parsed.Webserver != "" {
			asset.Props["webserver"] = parsed.Webserver
		}

		if parsed.CDNName != "" {
			asset.Props["cdn"] = parsed.CDNName
		}

		if parsed.Scheme != "" {
			asset.Props["scheme"] = parsed.Scheme
		}

		result = append(result, asset)

	}

	return result, invalid, nil

}
