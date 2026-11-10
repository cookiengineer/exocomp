package chat

import "context"
import "time"

// sleepRetry waits for the next attempt or returns early when the request
// context is cancelled.
func sleepRetry(ctx context.Context, delay time.Duration) error {

	if delay <= 0 {

		select {

		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil

		}

	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {

	case <-ctx.Done():
		return ctx.Err()

	case <-timer.C:
		return nil

	}

}
