package rand

import "strings"

func createStem(role string) string {

	role = strings.ToLower(strings.TrimSpace(role))
	role = strings.ReplaceAll(role, "-", "")
	role = strings.ReplaceAll(role, " ", "")

	for from, to := range compound_replacements {

		if role == from {
			role = to
			break
		}

	}

	for _, suffix := range unwanted_suffixes {

		if strings.HasSuffix(role, suffix) && len(role) > len(suffix)+2 {
			role = strings.TrimSuffix(role, suffix)
			break
		}

	}

	for from, to := range phonetic_stabilizers {

		if strings.Contains(role, from) {
			role = strings.ReplaceAll(role, from, to)
		}

	}

	switch {
	case strings.HasSuffix(role, "tt"):
		role = role[:len(role)-1]

	case strings.HasSuffix(role, "ss"):
		role = role[:len(role)-1]

	case strings.HasSuffix(role, "ck"):
		role = role[:len(role)-2] + "k"

	case strings.HasSuffix(role, "rv"):
		role += "a"

	case strings.HasSuffix(role, "rn"):
		role += "e"
	}

	if len(role) > 7 {
		role = role[:7]
	}

	if len(role) < 3 {
		role += "ar"
	}

	return role

}
