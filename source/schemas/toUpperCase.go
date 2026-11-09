package schemas

import "strings"

func toUpperCase(raw string) string {

	check := strings.ToLower(raw)

	// NOTE: Dirty Hack for stupid LLMs
	if strings.HasPrefix(check, "define") {

		tmp2 := strings.TrimSpace(raw[6:])

		return "Define" + strings.ToUpper(tmp2[0:1]) + strings.ToLower(tmp2[1:])

	} else {
		return strings.ToUpper(raw[0:1]) + strings.ToLower(raw[1:])
	}

}
