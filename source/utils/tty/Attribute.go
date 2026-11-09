package tty

type Attribute uint8

const (
	AttrBold Attribute = 1 << iota
	AttrDim
	AttrItalic
	AttrUnderline
	AttrReverse
	AttrStrikethrough
)
