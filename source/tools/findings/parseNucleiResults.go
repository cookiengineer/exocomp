package findings

import "encoding/json"
import "fmt"
import "strings"

func parseNucleiResults(data []byte, scan_id string) ([]nucleiImportResult, int, error) {

	trimmed := strings.TrimSpace(string(data))

	if trimmed == "" {
		return make([]nucleiImportResult, 0), 0, nil
	}

	lines := make([]json.RawMessage, 0)

	if strings.HasPrefix(trimmed, "[") == true {

		raw := make([]json.RawMessage, 0)

		if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
			return nil, 0, fmt.Errorf("Cannot parse nuclei results.")
		}

		lines = append(lines, raw...)

	} else {

		for _, line := range strings.Split(trimmed, "\n") {

			line = strings.TrimSpace(line)

			if line != "" {
				lines = append(lines, json.RawMessage(line))
			}

		}

	}

	result := make([]nucleiImportResult, 0, len(lines))
	invalid := 0

	for _, line := range lines {

		if strings.HasPrefix(string(line), "#") == true {
			continue
		}

		parsed := nucleiImportLine{}
		err := json.Unmarshal(line, &parsed)

		if err != nil {
			invalid++
			continue
		}

		if parsed.TemplateID == "" && parsed.Info.Name == "" && parsed.MatchedAt == "" && parsed.Host == "" {
			invalid++
			continue
		}

		result = append(result, nucleiImportLineToResult(parsed, scan_id))

	}

	return result, invalid, nil

}
