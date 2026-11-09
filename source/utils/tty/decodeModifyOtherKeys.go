package tty

import "strconv"
import "strings"

func decodeModifyOtherKeys(parameters string) (Key, error) {

	parts := strings.Split(parameters, ";")

	if len(parts) >= 3 {

		modifier, err1 := strconv.Atoi(parts[1])
		code, err2 := strconv.Atoi(parts[2])

		if err1 == nil && err2 == nil {

			shift, alt, ctrl := decodeModifier(modifier)

			switch code {

			case 9:
				return Key{Type: KeyTab, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
			case 13:
				return Key{Type: KeyEnter, Shift: shift, Alt: alt, Ctrl: ctrl}, nil
			default:
				return Key{Type: KeyRune, Rune: rune(code), Shift: shift, Alt: alt, Ctrl: ctrl}, nil

			}

		}

	}

	return Key{Type: KeyUnknown}, nil

}
