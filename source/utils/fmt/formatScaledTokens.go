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
