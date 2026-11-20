package assets

import "exocomp/types"
import "encoding/json"
import "strings"

func formatAssets(entries []types.Asset, format string) ([]byte, error) {

	if format == "json" {
		return json.MarshalIndent(entries, "", "\t")
	}

	lines := make([]string, 0, len(entries))

	for _, entry := range entries {

		value := entry.Value

		if format == "urls" {

			if entry.URL != "" {
				value = entry.URL
			} else if strings.HasPrefix(entry.Value, "http://") == false && strings.HasPrefix(entry.Value, "https://") == false {
				continue
			}

		}

		if strings.TrimSpace(value) != "" {
			lines = append(lines, value)
		}

	}

	return []byte(strings.Join(lines, "\n") + "\n"), nil

}
