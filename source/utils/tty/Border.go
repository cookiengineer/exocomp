package tty

type BorderChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	TopT        string
	BottomT     string
	LeftT       string
	RightT      string
	Cross       string
}

var BorderNone = BorderChars{
	Horizontal: " ",
}

var BorderLight = BorderChars{
	TopLeft:     "┌",
	TopRight:    "┐",
	BottomLeft:  "└",
	BottomRight: "┘",
	Horizontal:  "─",
	Vertical:    "│",
	TopT:        "┬",
	BottomT:     "┴",
	LeftT:       "├",
	RightT:      "┤",
	Cross:       "┼",
}

var BorderRounded = BorderChars{
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "╰",
	BottomRight: "╯",
	Horizontal:  "─",
	Vertical:    "│",
	TopT:        "┬",
	BottomT:     "┴",
	LeftT:       "├",
	RightT:      "┤",
	Cross:       "┼",
}

var BorderHeavy = BorderChars{
	TopLeft:     "┏",
	TopRight:    "┓",
	BottomLeft:  "┗",
	BottomRight: "┛",
	Horizontal:  "━",
	Vertical:    "┃",
	TopT:        "┳",
	BottomT:     "┻",
	LeftT:       "┣",
	RightT:      "┫",
	Cross:       "╋",
}

var BorderSplit = BorderChars{
	Horizontal: " ",
	Vertical:   "┃",
}

func DrawBox(buffer *Buffer, rect Rect, chars BorderChars, style Style) {

	if rect.Width < 2 || rect.Height < 2 {
		return
	}

	for x := rect.X + 1; x < rect.X+rect.Width-1; x++ {
		buffer.Set(x, rect.Y, chars.Horizontal, style)
		buffer.Set(x, rect.Y+rect.Height-1, chars.Horizontal, style)
	}

	for y := rect.Y + 1; y < rect.Y+rect.Height-1; y++ {
		buffer.Set(rect.X, y, chars.Vertical, style)
		buffer.Set(rect.X+rect.Width-1, y, chars.Vertical, style)
	}

	buffer.Set(rect.X, rect.Y, chars.TopLeft, style)
	buffer.Set(rect.X+rect.Width-1, rect.Y, chars.TopRight, style)
	buffer.Set(rect.X, rect.Y+rect.Height-1, chars.BottomLeft, style)
	buffer.Set(rect.X+rect.Width-1, rect.Y+rect.Height-1, chars.BottomRight, style)

}

func DrawSplit(buffer *Buffer, rect Rect, style Style) {

	for y := rect.Y; y < rect.Y+rect.Height; y++ {

		buffer.Set(rect.X, y, BorderSplit.Vertical, style)

		if rect.Width > 1 {
			buffer.Set(rect.X+rect.Width-1, y, BorderSplit.Vertical, style)
		}

	}

}

func DrawRule(buffer *Buffer, x int, y int, width int, style Style) {

	if width > 0 {
		buffer.WriteString(x, y, repeatText(BorderLight.Horizontal, width), style)
	}

}

func repeatText(text string, count int) string {

	result := make([]byte, 0)

	for count > 0 {
		result = append(result, text...)
		count--
	}

	return string(result)

}
