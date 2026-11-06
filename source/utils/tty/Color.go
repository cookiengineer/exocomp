package tty

import "fmt"
import "math"
import "os"
import "strconv"
import "strings"

type ColorProfile int

const (
	ProfileMono ColorProfile = iota
	Profile16
	Profile256
	ProfileTrueColor
)

type RGBA struct {
	R float64
	G float64
	B float64
	A float64
}

func NewRGBA(red float64, green float64, blue float64, alpha float64) RGBA {

	return RGBA{
		R: clamp(red, 0.0, 1.0),
		G: clamp(green, 0.0, 1.0),
		B: clamp(blue, 0.0, 1.0),
		A: clamp(alpha, 0.0, 1.0),
	}

}

func FromInts(red uint8, green uint8, blue uint8, alpha uint8) RGBA {

	return NewRGBA(
		float64(red)/255.0,
		float64(green)/255.0,
		float64(blue)/255.0,
		float64(alpha)/255.0,
	)

}

func FromHex(input string) (RGBA, error) {

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
		return RGBA{}, fmt.Errorf("tty.FromHex: Invalid color \"%s\".", input)
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
				return RGBA{}, fmt.Errorf("tty.FromHex: Invalid color \"%s\".", input)
			}

		}

		return FromInts(uint8(red), uint8(green), uint8(blue), uint8(alpha)), nil

	} else {
		return RGBA{}, fmt.Errorf("tty.FromHex: Invalid color \"%s\".", input)
	}

}

func (color RGBA) Hex() string {

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

func (color RGBA) WithAlpha(alpha float64) RGBA {
	return NewRGBA(color.R, color.G, color.B, alpha)
}

func (color RGBA) Opaque() bool {
	return color.A >= 1.0
}

func (color RGBA) Blend(background RGBA) RGBA {

	alpha := clamp(color.A, 0.0, 1.0)

	red := color.R*alpha + background.R*(1.0-alpha)
	green := color.G*alpha + background.G*(1.0-alpha)
	blue := color.B*alpha + background.B*(1.0-alpha)
	result_alpha := alpha + background.A*(1.0-alpha)

	return NewRGBA(red, green, blue, result_alpha)

}

func Tint(base RGBA, overlay RGBA, alpha float64) RGBA {

	amount := clamp(alpha, 0.0, 1.0)

	return NewRGBA(
		base.R+(overlay.R-base.R)*amount,
		base.G+(overlay.G-base.G)*amount,
		base.B+(overlay.B-base.B)*amount,
		base.A+(overlay.A-base.A)*amount,
	)

}

func DetectColorProfile() ColorProfile {

	no_color := strings.ToLower(strings.TrimSpace(os.Getenv("NO_COLOR")))

	if no_color == "1" || no_color == "true" || no_color == "yes" {
		return ProfileMono
	}

	color_term := strings.ToLower(strings.TrimSpace(os.Getenv("COLORTERM")))
	term := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))

	if color_term == "truecolor" || color_term == "24bit" || strings.Contains(term, "direct") {
		return ProfileTrueColor
	}

	if strings.Contains(term, "256color") {
		return Profile256
	}

	if term != "" && term != "dumb" {
		return Profile16
	}

	return ProfileMono

}

func (color RGBA) Foreground(profile ColorProfile) string {

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

func (color RGBA) Background(profile ColorProfile) string {

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

func clamp(value float64, minimum float64, maximum float64) float64 {

	if value < minimum {
		return minimum
	} else if value > maximum {
		return maximum
	} else {
		return value
	}

}
