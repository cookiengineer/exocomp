package fmt

import utils_strings "exocomp/utils/strings"
import "strings"

// FormatASCII transliterates non-ASCII characters to their ASCII equivalents
// while preserving tabs and newlines, so multi-line prompts keep their shape.
func FormatASCII(value string) string {

	tmp := strings.TrimSpace(value)
	filtered := ""

	for _, character := range tmp {

		if character == '\t' || character == '\n' {
			filtered += string(character)
		} else if mapped := utils_strings.Transliterate(character); mapped != "" {
			filtered += mapped
		}

	}

	return filtered

}
