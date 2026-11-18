package shell

import "testing"

func TestTerminalFilter(t *testing.T) {

	filter := &TerminalFilter{}

	result := filter.Write([]byte("hello\x1b[31mred\x1b[0m\r\nworld\x00"))

	if string(result) != "hellored\nworld" {
		t.Errorf("Expected filtered content, got %q", string(result))
	}

}
