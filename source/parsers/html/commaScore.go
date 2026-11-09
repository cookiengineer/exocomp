package html

import "strings"

func commaScore(text string) float64 {
	return float64(strings.Count(text, ","))
}
