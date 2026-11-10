package chat

import "context"
import "errors"
import "time"
import "testing"

func TestSleepRetry_CancelledContext(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	started := time.Now()
	err := sleepRetry(ctx, time.Second)

	if err != nil {

		if errors.Is(err, context.Canceled) == false {
			t.Errorf("Expected a cancelled context error, got %v", err)
		}

		if time.Since(started) > 100*time.Millisecond {
			t.Errorf("Expected cancellation to return immediately")
		}

	} else {
		t.Errorf("Expected an error")
	}

}

func TestSleepRetry_ZeroDelay(t *testing.T) {

	err := sleepRetry(context.Background(), 0)

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestSleepRetry_PositiveDelay(t *testing.T) {

	started := time.Now()
	err := sleepRetry(context.Background(), 5*time.Millisecond)

	if err == nil {

		if time.Since(started) < 5*time.Millisecond {
			t.Errorf("Expected sleepRetry to wait for the delay")
		}

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}
