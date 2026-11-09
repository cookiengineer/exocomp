package tty

import "os"
import "strings"

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
