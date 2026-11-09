package tty

import "time"

type Transition struct {
	Type           TransitionType
	Duration       time.Duration
	VisualDuration time.Duration
	Ease           Ease
}
