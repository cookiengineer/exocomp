package tty

import "time"

func Tween(duration time.Duration) Transition {

	return Transition{
		Type:     TransitionTween,
		Duration: duration,
		Ease:     Smoothstep,
	}

}
