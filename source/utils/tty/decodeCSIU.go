package tty

import "strconv"
import "strings"

func decodeCSIU(parameters string) (Key, error) {

	parts := strings.Split(parameters, ";")
	code, err := strconv.Atoi(parts[0])

	if err != nil {
		return Key{Type: KeyUnknown}, nil
	}

	modifier := 1

	if len(parts) > 1 {

		parsed, err1 := strconv.Atoi(parts[1])
		if err1 == nil {
			modifier = parsed
		}

	}

	shift, alt, ctrl := decodeModifier(modifier)

	switch code {

	case 9:
		return Key{Type: KeyTab, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 13:
		return Key{Type: KeyEnter, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 27:
		return Key{Type: KeyEscape, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	case 127:
		return Key{Type: KeyBackspace, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
	default:
		return Key{Type: KeyRune, Rune: rune(code), Shift: shift, Alt: alt, Ctrl: ctrl}, nil

	}

}
