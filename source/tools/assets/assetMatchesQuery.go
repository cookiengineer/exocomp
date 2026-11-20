package assets

import "exocomp/types"
import "strings"

func assetMatchesQuery(asset types.Asset, query string) bool {

	query = strings.ToLower(strings.TrimSpace(query))

	if query == "" {
		return true
	} else {

		values := []string{
			string(asset.Type),
			asset.Value,
			asset.Domain,
			asset.IP,
			asset.URL,
			asset.Endpoint,
			asset.Source,
		}

		values = append(values, asset.Technology...)
		values = append(values, asset.Tags...)

		for _, value := range asset.Props {
			values = append(values, value)
		}

		for _, value := range values {

			if strings.Contains(strings.ToLower(value), query) == true {
				return true
			}

		}

		return false

	}

}
