package tty

import "math"
import "sync"
import "time"

type Ease func(float64) float64

type TransitionType int

const (
	TransitionTween TransitionType = iota
	TransitionSpring
)

type Transition struct {
	Type           TransitionType
	Duration       time.Duration
	VisualDuration time.Duration
	Ease           Ease
}

func Tween(duration time.Duration) Transition {

	return Transition{
		Type:     TransitionTween,
		Duration: duration,
		Ease:     Smoothstep,
	}

}

func Spring(visual_duration time.Duration) Transition {

	return Transition{
		Type:           TransitionSpring,
		VisualDuration: visual_duration,
	}

}

func Smoothstep(progress float64) float64 {

	value := clamp(progress, 0.0, 1.0)

	return value * value * (3.0 - 2.0*value)

}

func Smootherstep(progress float64) float64 {

	value := clamp(progress, 0.0, 1.0)

	return value * value * value * (value*(value*6.0-15.0) + 10.0)

}

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

type Scheduler struct {
	mutex  *sync.Mutex
	stop   chan struct{}
	closed bool
}

func NewScheduler() *Scheduler {

	return &Scheduler{
		mutex: &sync.Mutex{},
		stop:  make(chan struct{}),
	}

}

func (scheduler *Scheduler) Start(interval time.Duration, callback func(time.Time)) {

	scheduler.mutex.Lock()

	if scheduler.closed == true {
		scheduler.mutex.Unlock()
		return
	}

	stop := scheduler.stop
	scheduler.mutex.Unlock()

	ticker := time.NewTicker(interval)

	go func() {

		defer ticker.Stop()

		for {

			select {
			case now := <-ticker.C:
				callback(now)
			case <-stop:
				return
			}

		}

	}()

}

func (scheduler *Scheduler) Stop() {

	scheduler.mutex.Lock()

	if scheduler.closed == false {
		close(scheduler.stop)
		scheduler.closed = true
	}

	scheduler.mutex.Unlock()

}
