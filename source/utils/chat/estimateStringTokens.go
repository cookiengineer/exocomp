package chat

import "strings"

func estimateStringTokens(str string, chars_per_token float64) int {

	if str != "" {

		tmp    := strings.TrimSpace(str)
		tokens := int(float64(len(tmp)) / chars_per_token)

		if tokens > 0 {
			return tokens
		} else {
			return 1
		}

	} else {
		return 0
	}

}
