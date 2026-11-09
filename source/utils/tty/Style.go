package tty

import "strings"

type Style struct {
	Foreground    Color
	Background    Color
	Attributes    Attribute
	HasForeground bool
	HasBackground bool
}

func NewStyle(foreground Color, background Color) Style {

	return Style{
		Foreground:    foreground,
		Background:    background,
		HasForeground: true,
		HasBackground: true,
	}

}

func (style Style) WithForeground(color Color) Style {

	style.Foreground = color
	style.HasForeground = true

	return style

}

func (style Style) WithBackground(color Color) Style {

	style.Background = color
	style.HasBackground = true

	return style

}

func (style Style) WithAttributes(attributes Attribute) Style {

	style.Attributes = style.Attributes | attributes

	return style

}

func (style Style) Reset() Style {

	return Style{}

}

func (style Style) Equals(other Style) bool {

	return style.Foreground == other.Foreground &&
		style.Background == other.Background &&
		style.Attributes == other.Attributes &&
		style.HasForeground == other.HasForeground &&
		style.HasBackground == other.HasBackground

}

func (style Style) Sequence(profile ColorProfile) string {

	if profile == ProfileMono {
		return ""
	}

	builder := strings.Builder{}

	builder.WriteString("\033[0m")

	if style.Attributes&AttrBold != 0 {
		builder.WriteString("\033[1m")
	}

	if style.Attributes&AttrDim != 0 {
		builder.WriteString("\033[2m")
	}

	if style.Attributes&AttrItalic != 0 {
		builder.WriteString("\033[3m")
	}

	if style.Attributes&AttrUnderline != 0 {
		builder.WriteString("\033[4m")
	}

	if style.Attributes&AttrReverse != 0 {
		builder.WriteString("\033[7m")
	}

	if style.Attributes&AttrStrikethrough != 0 {
		builder.WriteString("\033[9m")
	}

	if style.HasForeground == true {
		builder.WriteString(style.Foreground.Foreground(profile))
	}

	if style.HasBackground == true {
		builder.WriteString(style.Background.Background(profile))
	}

	return builder.String()

}
