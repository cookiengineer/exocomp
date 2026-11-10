package chat

import net_http "net/http"
import "math/rand/v2"
import "time"

// retryDelay returns the delay before the next attempt. A valid Retry-After
// header wins over the exponential backoff, but is capped to keep the agent
// responsive.
func retryDelay(attempt int, response *net_http.Response) time.Duration {

	if response != nil {

		delay, ok := retryAfterDelay(response.Header.Get("Retry-After"))

		if ok == true {

			if delay > RetryAfter {
				return RetryAfter
			} else {
				return delay
			}

		}

	}

	delay := 1 * time.Second

	for retry := 0; retry < attempt; retry++ {

		delay *= 2

		if delay >= RetryDelay {
			return RetryDelay
		}

	}

	if delay > 0 {
		return time.Duration(rand.Int64N(int64(delay) + 1))
	} else {
		return 0
	}

}
