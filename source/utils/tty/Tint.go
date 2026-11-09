package tty

func Tint(base Color, overlay Color, alpha float64) Color {

	amount := clamp(alpha, 0.0, 1.0)

	return NewColor(
		base.R+(overlay.R-base.R)*amount,
		base.G+(overlay.G-base.G)*amount,
		base.B+(overlay.B-base.B)*amount,
		base.A+(overlay.A-base.A)*amount,
	)

}
