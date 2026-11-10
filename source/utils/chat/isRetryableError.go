package chat

import "context"
import "errors"

// isRetryableError reports whether a transport error is worth retrying. A
// cancelled context is an explicit stop, so it is never retried.
func isRetryableError(err error) bool {

	if err == nil {
		return false
	} else if errors.Is(err, context.Canceled) == true {
		return false
	} else {
		return true
	}

}
