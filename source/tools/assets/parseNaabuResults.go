package assets

import "exocomp/types"
import utils_strings "exocomp/utils/strings"
import "encoding/json"
import "fmt"

func parseNaabuResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err == nil {

		result := make([]types.Asset, 0, len(lines))

		for _, line := range lines {

			parsed := naabuImportLine{}

			if json.Unmarshal(line, &parsed) != nil {
				invalid++
				continue
			}

			host := utils_strings.FirstNonEmpty(parsed.Host, parsed.IP)

			if host == "" || parsed.Port < 1 || parsed.Port > 65535 {
				invalid++
				continue
			}

			asset := newImportedAsset(types.AssetTypePort, fmt.Sprintf("%s:%d", host, parsed.Port), "naabu", in_scope)
			asset.Domain = host
			asset.IP = parsed.IP
			asset.Port = parsed.Port
			asset.Props = map[string]string{}

			if parsed.Protocol != "" {
				asset.Props["protocol"] = parsed.Protocol
			}

			if parsed.TLS == true {
				asset.Props["tls"] = "true"
			}

			result = append(result, asset)

		}

		return result, invalid, nil

	} else {
		return nil, 0, err
	}
}
