package chat

import net_http "net/http"
import "time"
import "testing"

func TestRetryAfterDelay_Seconds(t *testing.T) {

	result, ok := retryAfterDelay("2")

	if ok == true {

		if result != 2*time.Second {
			t.Errorf("Expected delay %s to be %s", result, 2*time.Second)
		}

	} else {
		t.Errorf("Expected seconds value to be valid")
	}

	result, ok = retryAfterDelay(" 2 ")

	if ok == true {

		if result != 2*time.Second {
			t.Errorf("Expected delay %s to be %s", result, 2*time.Second)
		}

	} else {
		t.Errorf("Expected trimmed seconds value to be valid")
	}

}

func TestRetryAfterDelay_HTTPDate(t *testing.T) {

	value := time.Now().Add(5 * time.Second).UTC().Format(net_http.TimeFormat)

	result, ok := retryAfterDelay(value)

	if ok == true {

		if result <= 0 || result > 5*time.Second {
			t.Errorf("Expected delay to be between 0 and %s, got %s", 5*time.Second, result)
		}

	} else {
		t.Fatalf("Expected date %q to be valid", value)
	}

}

func TestRetryAfterDelay_PastHTTPDate(t *testing.T) {

	value := time.Now().Add(-5 * time.Second).UTC().Format(net_http.TimeFormat)

	result, ok := retryAfterDelay(value)

	if ok == true {

		if result != 0 {
			t.Errorf("Expected past date delay %s to be %s", result, 0*time.Second)
		}

	} else {
		t.Errorf("Expected past date %q to be valid", value)
	}

}

func TestRetryAfterDelay_Invalid(t *testing.T) {

	invalid := []string{
		"",
		"tomorrow",
		"-1",
	}

	for _, value := range invalid {

		if _, ok := retryAfterDelay(value); ok == true {
			t.Errorf("Expected invalid Retry-After value %q to be rejected", value)
		}

	}

}
