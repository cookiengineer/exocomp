package fmt

import "testing"

func TestFormatASCII_Transliterates(t *testing.T) {

	result := FormatASCII("café crème")

	if result == "cafe creme" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "cafe creme", result)
	}

}

func TestFormatASCII_PreservesStructure(t *testing.T) {

	result := FormatASCII("first\n\tsecond")

	if result == "first\n\tsecond" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "first\n\tsecond", result)
	}

}

func TestFormatASCII_TrimsAndDropsUnknown(t *testing.T) {

	result := FormatASCII("  café  ")

	if result == "cafe" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "cafe", result)
	}

}

func TestFormatASCII_Empty(t *testing.T) {

	result := FormatASCII("")

	if result == "" {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", "", result)
	}

}
