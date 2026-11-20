package main

import "fmt"
import "go/token"
import "sort"

// Diagnostic is a single codestyle violation.
type Diagnostic struct {
	Position token.Position
	Rule     string
	Message  string
}

func (diagnostic Diagnostic) String() string {
	return fmt.Sprintf(
		"%s:%d:%d: [%s] %s",
		diagnostic.Position.Filename,
		diagnostic.Position.Line,
		diagnostic.Position.Column,
		diagnostic.Rule,
		diagnostic.Message,
	)
}

func sortDiagnostics(diagnostics []Diagnostic) {

	sort.SliceStable(diagnostics, func(a int, b int) bool {

		left := diagnostics[a]
		right := diagnostics[b]

		if left.Position.Filename != right.Position.Filename {
			return left.Position.Filename < right.Position.Filename
		} else if left.Position.Line != right.Position.Line {
			return left.Position.Line < right.Position.Line
		} else if left.Position.Column != right.Position.Column {
			return left.Position.Column < right.Position.Column
		} else {
			return left.Rule < right.Rule
		}

	})

}
