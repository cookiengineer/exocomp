package fmt

import "fmt"

func FormatPercentage(value float64, limit float64) string {

	percentage := float64(0.0)

	if limit > 0 {

		percentage = float64(value) / float64(limit)

		if percentage < 0.0 {
			percentage = 0.0
		} else if percentage > 1.0 {
			percentage = 1.0
		}

	}

	return fmt.Sprintf("%3.0f%%", percentage*100)

}
