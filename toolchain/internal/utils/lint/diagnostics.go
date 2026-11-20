package lint

import "fmt"
import "go/token"

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
