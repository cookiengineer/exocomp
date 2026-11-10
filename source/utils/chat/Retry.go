package chat

import "time"

// The complete call, including all retries and their delays, is bounded by a
// 10 minute timeout so a stalled connection cannot block an agent forever.
const (
	RetryMaxAttempts    = 5
	RetryDelay          = 30 * time.Second
	RetryAfter          = 60 * time.Second
	RetryRequestTimeout = 10 * time.Minute
)
