package tty

import "sync"
import "time"

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
