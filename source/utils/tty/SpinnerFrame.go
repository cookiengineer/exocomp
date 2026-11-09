package tty

import "time"

func SpinnerFrame(elapsed time.Duration, interval time.Duration) string {

	if interval <= 0 {
		interval = 80 * time.Millisecond
	}

	index := int(elapsed/interval) % len(SpinnerFrames)

	return SpinnerFrames[index]

}
