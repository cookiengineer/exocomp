package skill

import "exocomp/types"
import "strings"

func cleanYAMLString(input string) string {

	result := strings.TrimSpace(input)

	if strings.HasPrefix(result, "\"") && strings.HasSuffix(result, "\"") {
		result = strings.TrimSpace(result[1:len(result)-1])
	}

	if strings.HasPrefix(result, "'") && strings.HasSuffix(result, "'") {
		result = strings.TrimSpace(result[1:len(result)-1])
	}

	return strings.TrimSpace(result)

}

func ParseSkill(buffer []byte) *types.Skill {

	lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")

	if len(lines) > 2 {

		if strings.TrimSpace(lines[0]) == "---" {

			is_header := true
			header    := make([]string, 0)
			body      := make([]string, 0)

			for l := 1; l < len(lines); l++ {

				line := lines[l]

				if is_header == true {

					if strings.TrimSpace(line) == "---" {
						is_header = false
					} else {
						header = append(header, line)
					}

				} else if is_header == false {
					body = append(body, line)
				}

			}

			name          := ""
			description   := ""
			license       := ""
			compatibility := ""
			metadata      := make(map[string]string)
			allowed_tools := make([]string, 0)

			is_metadata := false

			for _, line := range header {

				if is_metadata == true {

					if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {

						tmp := strings.Split(strings.TrimSpace(line), ":")

						if len(tmp) >= 2 {

							key := strings.TrimSpace(tmp[0])
							val := strings.TrimSpace(strings.Join(tmp[1:], ":"))

							metadata[key] = val

						}

					}

				} else if strings.HasPrefix(line, "name: ") {

					name        = cleanYAMLString(line[6:])
					is_metadata = false

				} else if strings.HasPrefix(line, "description: ") {

					description = cleanYAMLString(line[13:])
					is_metadata = false

				} else if strings.HasPrefix(line, "license: ") {

					license     = cleanYAMLString(line[9:])
					is_metadata = false

				} else if strings.HasPrefix(line, "compatibility: ") {

					compatibility = cleanYAMLString(line[15:])
					is_metadata   = false

				} else if strings.HasPrefix(line, "metadata: ") {

					is_metadata = true

				} else if strings.HasPrefix(line, "allowed-tools: ") {

					tmp := strings.Split(line[15:], " ")

					for _, raw := range tmp {

						tool := strings.TrimSpace(raw)

						if strings.Contains(tool, ".") {
							allowed_tools = append(allowed_tools, tool)
						}

					}

					is_metadata = false

				} else if strings.TrimSpace(line) == "" {
					is_metadata = false
				}

			}

			return &types.Skill{
				Name:          name,
				Description:   description,
				License:       license,
				Compatibility: compatibility,
				Metadata:      metadata,
				AllowedTools:  allowed_tools,
				Body:          strings.Join(body, "\n"),
				Scripts:       make(map[string]string),
			}

		} else {
			return nil
		}

	} else {
		return nil
	}



}
