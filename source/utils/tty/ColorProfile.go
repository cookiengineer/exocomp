package tty

type ColorProfile int

const (
	ProfileMono ColorProfile = iota
	Profile16
	Profile256
	ProfileTrueColor
)
