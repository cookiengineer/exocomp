package chat

import net_http "net/http"
import "strconv"
import "strings"
import "time"

// retryAfterDelay parses a Retry-After header in either delta-seconds or
// HTTP-date form.
func retryAfterDelay(value string) (time.Duration, bool) {

	value = strings.TrimSpace(value)

	if value != "" {

		seconds, err0 := strconv.Atoi(value)

		if err0 == nil {

			if seconds < 0 {
				return 0, false
			} else {
				return time.Duration(seconds) * time.Second, true
			}

		}

		timestamp, err1 := net_http.ParseTime(value)

		if err1 == nil {

			delay := time.Until(timestamp)

			if delay < 0 {
				delay = 0
			}

			return delay, true

		}

	}

	return 0, false

}
