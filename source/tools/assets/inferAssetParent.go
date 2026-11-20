package assets

import "exocomp/types"
import "strings"

func inferAssetParent(tool *Assets, asset types.Asset) string {

	if asset.Domain != "" {

		domain := strings.ToLower(strings.TrimSpace(asset.Domain))

		for _, candidate := range tool.contents {

			if candidate.Type != types.AssetTypeDomain && candidate.Type != types.AssetTypeSubdomain && candidate.Type != types.AssetTypeDNSRecord {
				continue
			} else if strings.EqualFold(candidate.Value, domain) == true {
				return candidate.ID
			}
		}

		parts := strings.Split(domain, ".")

		if len(parts) > 2 {

			parent := strings.Join(parts[1:], ".")

			for _, candidate := range tool.contents {

				if candidate.Type != types.AssetTypeDomain && candidate.Type != types.AssetTypeSubdomain && candidate.Type != types.AssetTypeDNSRecord {
					continue
				}

				if strings.EqualFold(candidate.Value, parent) == true {
					return candidate.ID
				}

			}

		}

	}

	if asset.IP != "" {

		for _, candidate := range tool.contents {

			if candidate.Type == types.AssetTypeIP && strings.EqualFold(candidate.Value, asset.IP) == true {
				return candidate.ID
			}

		}

	}

	return ""

}
