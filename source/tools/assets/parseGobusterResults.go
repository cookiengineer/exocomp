package assets

import "exocomp/types"
import "strings"

func parseGobusterResults(data []byte, target string, in_scope bool) ([]types.Asset, int, error) {

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	result := make([]types.Asset, 0, len(lines))
	invalid := 0

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		marker := strings.Index(line, " (Status:")

		if marker < 0 {
			continue
		}

		value := strings.TrimSpace(line[:marker])
		status := ""

		if relative_end := strings.Index(line[marker:], ")"); relative_end >= 0 {
			status = strings.TrimSpace(line[marker+len(" (Status:") : marker+relative_end])
		}

		if value == "" {
			invalid++
			continue
		}

		if strings.HasPrefix(value, "http://") == false && strings.HasPrefix(value, "https://") == false {

			if target != "" {
				value = strings.TrimRight(target, "/") + "/" + strings.TrimLeft(value, "/")
			}

		}

		asset := newImportedAsset(types.AssetTypeEndpoint, value, "gobuster", in_scope)
		asset.URL = value
		asset.Endpoint = value
		asset.Domain = hostFromValue(value)
		asset.Props = map[string]string{}

		if status != "" {
			asset.Props["status_code"] = status
		}

		if size_marker := strings.Index(line, "[Size:"); size_marker >= 0 {

			if end := strings.Index(line[size_marker:], "]"); end >= 0 {
				size := strings.TrimSpace(line[size_marker+len("[Size:") : size_marker+end])
				asset.Props["size"] = size
			}

		}

		result = append(result, asset)

	}

	return result, invalid, nil

}
