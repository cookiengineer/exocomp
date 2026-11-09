package tty

import "fmt"
import "math"
import "strconv"
import "strings"

type Color struct {
	R float64
	G float64
	B float64
	A float64
}

func NewColor(red float64, green float64, blue float64, alpha float64) Color {

	return Color{
		R: clamp(red, 0.0, 1.0),
		G: clamp(green, 0.0, 1.0),
		B: clamp(blue, 0.0, 1.0),
		A: clamp(alpha, 0.0, 1.0),
	}

}

func ColorFromInts(red uint8, green uint8, blue uint8, alpha uint8) Color {

	return NewColor(
		float64(red)/255.0,
		float64(green)/255.0,
		float64(blue)/255.0,
		float64(alpha)/255.0,
	)

}

func ColorFromHex(input string) (Color, error) {

	value := strings.TrimSpace(input)

	if strings.HasPrefix(value, "#") {
		value = value[1:]
	}

	if len(value) == 3 {
		value = string([]byte{
			value[0], value[0],
			value[1], value[1],
			value[2], value[2],
		})
	} else if len(value) == 4 {
		value = string([]byte{
			value[0], value[0],
			value[1], value[1],
			value[2], value[2],
			value[3], value[3],
		})
	}

	if len(value) != 6 && len(value) != 8 {
		return Color{}, fmt.Errorf("tty.ColorFromHex: Invalid color \"%s\".", input)
	}

	red, err1 := strconv.ParseUint(value[0:2], 16, 8)
	green, err2 := strconv.ParseUint(value[2:4], 16, 8)
	blue, err3 := strconv.ParseUint(value[4:6], 16, 8)

	if err1 == nil && err2 == nil && err3 == nil {

		alpha := uint64(255)

		if len(value) == 8 {

			parsed, err4 := strconv.ParseUint(value[6:8], 16, 8)

			if err4 == nil {
				alpha = parsed
			} else {
				return Color{}, fmt.Errorf("tty.ColorFromHex: Invalid color \"%s\".", input)
			}

		}

		return ColorFromInts(uint8(red), uint8(green), uint8(blue), uint8(alpha)), nil

	} else {
		return Color{}, fmt.Errorf("tty.ColorFromHex: Invalid color \"%s\".", input)
	}

}

func (color Color) Hex() string {

	red := int(math.Round(color.R * 255.0))
	green := int(math.Round(color.G * 255.0))
	blue := int(math.Round(color.B * 255.0))

	if color.A < 1.0 {

		alpha := int(math.Round(color.A * 255.0))

		return fmt.Sprintf("#%02x%02x%02x%02x", red, green, blue, alpha)

	} else {
		return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
	}

}

func (color Color) WithAlpha(alpha float64) Color {
	return NewColor(color.R, color.G, color.B, alpha)
}

func (color Color) Opaque() bool {
	return color.A >= 1.0
}

func (color Color) Blend(background Color) Color {

	alpha := clamp(color.A, 0.0, 1.0)

	red := color.R*alpha + background.R*(1.0-alpha)
	green := color.G*alpha + background.G*(1.0-alpha)
	blue := color.B*alpha + background.B*(1.0-alpha)
	result_alpha := alpha + background.A*(1.0-alpha)

	return NewColor(red, green, blue, result_alpha)

}

func (color Color) Foreground(profile ColorProfile) string {

	if profile == ProfileMono {
		return ""
	}

	red := int(math.Round(color.R * 255.0))
	green := int(math.Round(color.G * 255.0))
	blue := int(math.Round(color.B * 255.0))

	if profile == ProfileTrueColor {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", red, green, blue)
	}

	code := nearestANSI(red, green, blue, profile)

	if code < 16 {
		return fmt.Sprintf("\033[%dm", 30+code)
	}

	return fmt.Sprintf("\033[38;5;%dm", code)

}

func (color Color) Background(profile ColorProfile) string {

	if profile == ProfileMono {
		return ""
	}

	red := int(math.Round(color.R * 255.0))
	green := int(math.Round(color.G * 255.0))
	blue := int(math.Round(color.B * 255.0))

	if profile == ProfileTrueColor {
		return fmt.Sprintf("\033[48;2;%d;%d;%dm", red, green, blue)
	}

	code := nearestANSI(red, green, blue, profile)

	if code < 16 {
		return fmt.Sprintf("\033[%dm", 40+code)
	}

	return fmt.Sprintf("\033[48;5;%dm", code)

}
