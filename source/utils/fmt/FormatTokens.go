package fmt

import "fmt"
import "strings"

func formatScaledTokens(value float64, suffix string) string {

	if value >= 10.0 {

		return fmt.Sprintf("%.0f%s", value, suffix)

	} else {

		text := fmt.Sprintf("%.1f", value)
		text = strings.TrimSuffix(text, ".0")

		return fmt.Sprintf("%s%s", text, suffix)

	}

}

func FormatTokens(value int) string {

	if value < 1000 {
		return fmt.Sprintf("%d", value)
	} else if value < 1000000 {
		return formatScaledTokens(float64(value)/1000.0, "k")
	} else if value < 1000000000 {
		return formatScaledTokens(float64(value)/1000000.0, "M")
	} else {
		return formatScaledTokens(float64(value)/1000000000.0, "G")
	}

}

