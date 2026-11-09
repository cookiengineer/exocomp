package tty

import "time"

func Spring(visual_duration time.Duration) Transition {

	return Transition{
		Type:           TransitionSpring,
		VisualDuration: visual_duration,
	}

}
