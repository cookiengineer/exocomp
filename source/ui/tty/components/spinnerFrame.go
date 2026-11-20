package components

import utils_tty "exocomp/utils/tty"
import "time"

func spinnerFrame(elapsed time.Duration) string {

	if len(utils_tty.SpinnerFrames) == 0 {
		return " "
	} else {

		index := int(elapsed/(80*time.Millisecond)) % len(utils_tty.SpinnerFrames)

		if index < 0 {
			index = 0
		}

		return utils_tty.SpinnerFrames[index]

	}
}
