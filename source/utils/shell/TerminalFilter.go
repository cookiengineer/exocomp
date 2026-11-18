package shell

import "strings"
import "unicode/utf8"

// TerminalFilter removes ANSI escape sequences and non-printable control
// characters from a PTY stream. It is stateful so escape sequences split
// across read chunks are still removed.
type TerminalFilter struct {
	state     int
	pendingCR bool
}

const (
	terminalStateNormal = iota
	terminalStateEscape
	terminalStateCSI
	terminalStateOSC
	terminalStateOSCEscape
	terminalStateString
	terminalStateStringEscape
)

// Write filters a chunk of terminal output.
func (filter *TerminalFilter) Write(data []byte) []byte {

	result := make([]byte, 0, len(data))

	for _, value := range data {

		if filter.pendingCR == true {

			filter.pendingCR = false

			if value == '\n' {
				result = append(result, '\n')
				continue
			} else {
				result = append(result, '\n')
			}

		}

		switch filter.state {

		case terminalStateNormal:

			if value == 0x1b {
				filter.state = terminalStateEscape
			} else if value == '\r' {
				filter.pendingCR = true
			} else if value == '\n' || value == '\t' {
				result = append(result, value)
			} else if value >= 0x20 && value != 0x7f {
				result = append(result, value)
			}

		case terminalStateEscape:

			if value == '[' {
				filter.state = terminalStateCSI
			} else if value == ']' {
				filter.state = terminalStateOSC
			} else if value == 'P' || value == '^' || value == '_' {
				filter.state = terminalStateString
			} else {
				filter.state = terminalStateNormal
			}

		case terminalStateCSI:

			if value >= 0x40 && value <= 0x7e {
				filter.state = terminalStateNormal
			}

		case terminalStateOSC:

			if value == 0x07 {
				filter.state = terminalStateNormal
			} else if value == 0x1b {
				filter.state = terminalStateOSCEscape
			}

		case terminalStateOSCEscape:

			if value == '\\' {
				filter.state = terminalStateNormal
			} else {
				filter.state = terminalStateOSC
			}

		case terminalStateString:

			if value == 0x1b {
				filter.state = terminalStateStringEscape
			}

		case terminalStateStringEscape:

			if value == '\\' {
				filter.state = terminalStateNormal
			} else {
				filter.state = terminalStateString
			}

		}

	}

	if filter.pendingCR == true {
		filter.pendingCR = false
		result = append(result, '\n')
	}

	return []byte(strings.ToValidUTF8(string(result), string(utf8.RuneError)))

}
