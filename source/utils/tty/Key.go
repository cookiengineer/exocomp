package tty

type KeyType int

const (
	KeyUnknown KeyType = iota
	KeyRune
	KeyEnter
	KeyTab
	KeyBackspace
	KeyDelete
	KeyEscape
	KeyInsert
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyPaste
	KeyEOF
)

type Key struct {
	Type  KeyType
	Rune  rune
	Text  string
	Ctrl  bool
	Alt   bool
	Shift bool
}

func (key Key) IsCtrl(value rune) bool {
	return key.Ctrl == true && key.Rune == value
}
