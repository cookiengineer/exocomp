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
	KeyMouse
	KeyPaste
	KeyEOF
)
