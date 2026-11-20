package tty

import "strconv"
import "strings"

func decodeSGRMouse(parameters string, final byte) (Key, error) {

	if strings.HasPrefix(parameters, "<") == false {
		return Key{Type: KeyUnknown}, nil
	} else {

		parts := strings.Split(parameters[1:], ";")

		if len(parts) != 3 {
			return Key{Type: KeyUnknown}, nil
		} else {

			button, err0 := strconv.Atoi(parts[0])
			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])

			if err0 != nil || err1 != nil || err2 != nil {
				return Key{Type: KeyUnknown}, nil
			} else {

				mouse_button := button & 3
				mouse_wheel := 0

				if button&64 != 0 {

					mouse_button = -1

					if button&1 == 0 {
						mouse_wheel = -1
					} else {
						mouse_wheel = 1
					}

				}

				return Key{
					Type:         KeyMouse,
					MouseButton:  mouse_button,
					MouseX:       max(0, x-1),
					MouseY:       max(0, y-1),
					MouseRelease: final == 'm' || button == 3,
					MouseWheel:   mouse_wheel,
				}, nil

			}

		}

	}

}
