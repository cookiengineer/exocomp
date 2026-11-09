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
