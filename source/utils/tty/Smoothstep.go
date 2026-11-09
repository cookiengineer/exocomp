package tty

func Smoothstep(progress float64) float64 {

	value := clamp(progress, 0.0, 1.0)

	return value * value * (3.0 - 2.0*value)

}
