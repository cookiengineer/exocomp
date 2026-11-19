package tty

type Key struct {
	Type         KeyType
	Rune         rune
	Text         string
	Ctrl         bool
	Alt          bool
	Shift        bool
	MouseButton  int
	MouseX       int
	MouseY       int
	MouseRelease bool
	MouseWheel   int
}

func (key Key) IsCtrl(value rune) bool {
	return key.Ctrl == true && key.Rune == value
}
