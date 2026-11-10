package chat

import net_http "net/http"
import "time"
import "testing"

func TestRetryDelay_HonorsRetryAfter(t *testing.T) {

	response := &net_http.Response{
		Header: net_http.Header{"Retry-After": []string{"2"}},
	}

	if result := retryDelay(0, response); result != 2*time.Second {
		t.Errorf("Expected delay %s to be %s", result, 2*time.Second)
	}

	response.Header.Set("Retry-After", "9999")

	if result := retryDelay(0, response); result != RetryAfter {
		t.Errorf("Expected delay %s to be %s", result, RetryAfter)
	}

	response.Header.Set("Retry-After", "0")

	if result := retryDelay(0, response); result != 0 {
		t.Errorf("Expected delay %s to be %s", result, 0*time.Second)
	}

}

func TestRetryDelay_ExponentialJitter(t *testing.T) {

	for attempt := 0; attempt < RetryMaxAttempts; attempt++ {

		result := retryDelay(attempt, nil)
		maximum := time.Second << attempt

		if maximum > RetryDelay {
			maximum = RetryDelay
		}

		if result < 0 || result > maximum {
			t.Errorf("Expected delay for attempt %d to be between 0 and %s, got %s", attempt, maximum, result)
		}

	}

	if result := retryDelay(RetryMaxAttempts, nil); result != RetryDelay {
		t.Errorf("Expected delay for attempt %d to be %s, got %s", RetryMaxAttempts, RetryDelay, result)
	}

	if result := retryDelay(RetryMaxAttempts+1, nil); result != RetryDelay {
		t.Errorf("Expected delay for attempt %d to be %s, got %s", RetryMaxAttempts+1, RetryDelay, result)
	}

}
