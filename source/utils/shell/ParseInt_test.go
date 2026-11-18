package shell

import "testing"

func TestParseInt(t *testing.T) {

	arguments := map[string]any{"count": float64(3)}

	result, err0 := ParseInt(arguments, "count", 1)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}

	fallback, err1 := ParseInt(arguments, "missing", 7)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if fallback != 7 {
		t.Errorf("Expected 7, got %d", fallback)
	}

	_, err2 := ParseInt(map[string]any{"count": 1.5}, "count", 1)

	if err2 == nil {
		t.Errorf("Expected fractional integer to fail")
	}

}
