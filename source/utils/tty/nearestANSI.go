package tty

import "math"

func nearestANSI(red int, green int, blue int, profile ColorProfile) int {

	if profile == Profile16 {

		palette := [16][3]int{
			{0, 0, 0},
			{128, 0, 0},
			{0, 128, 0},
			{128, 128, 0},
			{0, 0, 128},
			{128, 0, 128},
			{0, 128, 128},
			{192, 192, 192},
			{128, 128, 128},
			{255, 0, 0},
			{0, 255, 0},
			{255, 255, 0},
			{0, 0, 255},
			{255, 0, 255},
			{0, 255, 255},
			{255, 255, 255},
		}

		best_index := 0
		best_distance := math.MaxFloat64

		for index, entry := range palette {

			distance := float64((red-entry[0])*(red-entry[0]) + (green-entry[1])*(green-entry[1]) + (blue-entry[2])*(blue-entry[2]))

			if distance < best_distance {

				best_distance = distance
				best_index = index

			}

		}

		return best_index

	}

	if red == green && green == blue {

		if red < 8 {
			return 16
		} else if red > 248 {
			return 231
		} else {
			return 232 + ((red - 8) / 10)
		}

	}

	to_cube := func(value int) int {

		if value < 48 {
			return 0
		} else if value < 115 {
			return 1
		} else {
			return (value - 35) / 40
		}

	}

	ri := to_cube(red)
	gi := to_cube(green)
	bi := to_cube(blue)

	return 16 + 36*ri + 6*gi + bi

}
