package tty

func Smootherstep(progress float64) float64 {

	value := clamp(progress, 0.0, 1.0)

	return value * value * value * (value*(value*6.0-15.0) + 10.0)

}
