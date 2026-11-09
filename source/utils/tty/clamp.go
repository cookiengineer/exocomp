package tty

func clamp(value float64, minimum float64, maximum float64) float64 {

	if value < minimum {
		return minimum
	} else if value > maximum {
		return maximum
	} else {
		return value
	}

}
