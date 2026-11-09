package tty

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
