package fmt

import "fmt"

func FormatUSD(amount float64) string {

	if amount < 0.01 {
		return fmt.Sprintf("$%.4f", amount)
	} else {
		return fmt.Sprintf("$%.2f", amount)
	}

}
