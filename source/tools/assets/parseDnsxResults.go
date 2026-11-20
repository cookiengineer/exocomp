package assets

import "exocomp/types"
import "encoding/json"
import "strings"

func parseDnsxResults(data []byte, in_scope bool) ([]types.Asset, int, error) {

	lines, invalid, err := decodeImportLines(data)

	if err != nil {
		return nil, 0, err
	}

	result := make([]types.Asset, 0, len(lines))

	for _, line := range lines {

		parsed := dnsxImportLine{}

		if json.Unmarshal(line, &parsed) != nil {
			invalid++
			continue
		}

		if parsed.Host == "" {
			invalid++
			continue
		}

		asset := newImportedAsset(types.AssetTypeDNSRecord, parsed.Host, "dnsx", in_scope)
		asset.Domain = parsed.Host
		asset.Props = map[string]string{}

		if len(parsed.A) > 0 {
			asset.Props["a"] = strings.Join(parsed.A, ", ")
		}

		if len(parsed.AAAA) > 0 {
			asset.Props["aaaa"] = strings.Join(parsed.AAAA, ", ")
		}

		if len(parsed.CNAME) > 0 {
			asset.Props["cname"] = strings.Join(parsed.CNAME, ", ")
		}

		if len(parsed.MX) > 0 {
			asset.Props["mx"] = strings.Join(parsed.MX, ", ")
		}

		if len(parsed.NS) > 0 {
			asset.Props["ns"] = strings.Join(parsed.NS, ", ")
		}

		if len(parsed.TXT) > 0 {
			asset.Props["txt"] = strings.Join(parsed.TXT, ", ")
		}

		if parsed.StatusCode != "" {
			asset.Props["status_code"] = parsed.StatusCode
		}

		result = append(result, asset)

		for _, ip := range parsed.A {

			ip_asset := newImportedAsset(types.AssetTypeIP, ip, "dnsx", in_scope)
			ip_asset.Domain = parsed.Host
			result = append(result, ip_asset)

		}

		for _, ip := range parsed.AAAA {

			ip_asset := newImportedAsset(types.AssetTypeIP, ip, "dnsx", in_scope)
			ip_asset.Domain = parsed.Host
			result = append(result, ip_asset)

		}

	}

	return result, invalid, nil

}
