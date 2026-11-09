package fmt

import "fmt"

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
