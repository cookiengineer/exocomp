package findings

import "strings"
import "testing"
import "unicode/utf8"

func TestLimitText(t *testing.T) {

	short := LimitText(" short ", 16)

	if short != "short" {
		t.Errorf("Expected trimmed short text, got %q", short)
	}

	long := LimitText(strings.Repeat("a", 32), 16)

	if strings.HasPrefix(long, strings.Repeat("a", 16)) == false {
		t.Errorf("Expected truncated text, got %q", long)
	}

	if strings.HasSuffix(long, "[TRUNCATED]") == false {
		t.Errorf("Expected truncation marker, got %q", long)
	}

	unicode_text := LimitText(strings.Repeat("ä", 20), 16)

	if utf8.ValidString(unicode_text) == false {
		t.Errorf("Expected valid UTF-8 after truncation")
	}

}
