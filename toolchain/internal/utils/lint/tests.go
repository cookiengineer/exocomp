package lint

import "regexp"

var testFunctionPattern = regexp.MustCompile(`^Test[A-Z][A-Za-z0-9]*(_[A-Za-z0-9]+)*$`)
