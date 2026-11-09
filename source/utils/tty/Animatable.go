package tty

import "math"
import "sync"
import "time"

type Animatable struct {
	mutex      *sync.Mutex
	transition Transition
	value      float64
	target     float64
	from       float64
	velocity   float64
	started    time.Time
	previous   time.Time
	active     bool
}

func NewAnimatable(initial float64, transition Transition) *Animatable {

	now := time.Now()

	return &Animatable{
		mutex:      &sync.Mutex{},
		transition: transition,
		value:      initial,
		target:     initial,
		from:       initial,
		started:    now,
		previous:   now,
	}

}

func (animatable *Animatable) Value() float64 {

	animatable.mutex.Lock()
	value := animatable.value
	animatable.mutex.Unlock()

	return value

}

func (animatable *Animatable) Active() bool {

	animatable.mutex.Lock()
	active := animatable.active
	animatable.mutex.Unlock()

	return active

}

func (animatable *Animatable) Jump(value float64) {

	animatable.mutex.Lock()
	animatable.value = value
	animatable.target = value
	animatable.from = value
	animatable.velocity = 0.0
	animatable.active = false
	animatable.mutex.Unlock()

}

func (animatable *Animatable) Animate(target float64) {

	animatable.mutex.Lock()

	if target == animatable.target {

		animatable.mutex.Unlock()
		return

	}

	animatable.target = target
	animatable.from = animatable.value
	animatable.started = time.Now()
	animatable.previous = animatable.started
	animatable.active = true
	animatable.mutex.Unlock()

}

func (animatable *Animatable) Update(now time.Time) {

	animatable.mutex.Lock()
	defer animatable.mutex.Unlock()

	if animatable.active == false {
		return
	}

	if animatable.transition.Type == TransitionTween {

		duration := animatable.transition.Duration

		if duration <= 0 {

			animatable.value = animatable.target
			animatable.active = false
			return

		}

		progress := float64(now.Sub(animatable.started)) / float64(duration)
		eased := animatable.transition.Ease(progress)

		animatable.value = animatable.from + (animatable.target-animatable.from)*eased

		if progress >= 1.0 {

			animatable.value = animatable.target
			animatable.active = false

		}

	} else {

		delta := math.Min(0.05, now.Sub(animatable.previous).Seconds())
		animatable.previous = now

		duration := animatable.transition.VisualDuration

		if duration <= 0 {
			duration = 160 * time.Millisecond
		}

		frequency := (2.0 * math.Pi) / (math.Max(0.001, duration.Seconds()) * 1.2)
		offset := animatable.value - animatable.target
		decay := math.Exp(-frequency * delta)
		next_velocity := animatable.velocity + frequency*offset

		animatable.value = animatable.target + (offset+next_velocity*delta)*decay
		animatable.velocity = (animatable.velocity - frequency*next_velocity*delta) * decay

		if math.Abs(animatable.value-animatable.target) < 0.002 && math.Abs(animatable.velocity) < 0.002 {

			animatable.value = animatable.target
			animatable.velocity = 0.0
			animatable.active = false

		}

	}

}
