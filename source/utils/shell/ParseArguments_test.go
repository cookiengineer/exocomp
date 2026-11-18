package shell

import "testing"

func TestParseArguments(t *testing.T) {

	arguments, err0 := ParseArguments([]any{"one", "two"})

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	if len(arguments) != 2 || arguments[0] != "one" || arguments[1] != "two" {
		t.Errorf("Expected [one two], got %v", arguments)
	}

	_, err1 := ParseArguments([]any{"one", 2})

	if err1 == nil {
		t.Errorf("Expected invalid argument to fail")
	}

}

func TestArgumentBytes(t *testing.T) {

	result := ArgumentBytes([]string{"abc", "de"})

	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}

}
