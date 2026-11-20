package strings

import "testing"

func TestToASCII_Transliterates(t *testing.T) {

	result := ToASCII("café crème")

	if result == "cafe creme" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "cafe creme", result)
	}

}

func TestToASCII_TabBecomesSpaces(t *testing.T) {

	result := ToASCII("a\tb")

	if result == "a    b" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "a    b", result)
	}

}

func TestToASCII_DropsNewlines(t *testing.T) {

	result := ToASCII("a\nb")

	if result == "ab" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "ab", result)
	}

}

func TestToASCII_Empty(t *testing.T) {

	result := ToASCII("")

	if result == "" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "", result)
	}

}
