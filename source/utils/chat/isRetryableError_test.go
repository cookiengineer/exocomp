package chat

import "context"
import "errors"
import "fmt"
import "testing"

func TestIsRetryableError(t *testing.T) {

	if isRetryableError(nil) != false {
		t.Errorf("Expected %v to not be retryable", nil)
	}

	if isRetryableError(context.Canceled) != false {
		t.Errorf("Expected cancelled context to not be retryable")
	}

	if isRetryableError(fmt.Errorf("request failed: %w", context.Canceled)) != false {
		t.Errorf("Expected wrapped cancelled context to not be retryable")
	}

	if isRetryableError(errors.New("connection reset")) != true {
		t.Errorf("Expected a transport error to be retryable")
	}

}
