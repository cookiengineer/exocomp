package lint

import "regexp"

var initialismPattern = regexp.MustCompile(`^[A-Z]{2,5}$`)

var allowedInitialisms = map[string]bool{
	"ANSI":  true,
	"API":   true,
	"ASCII": true,
	"CIDR":  true,
	"CPU":   true,
	"HTML":  true,
	"HTTP":  true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"OS":    true,
	"TTY":   true,
	"UI":    true,
	"URL":   true,
	"YAML":  true,
}

var predicatePrefixes = []string{
	"allow_", "can_", "did_", "has_", "in_", "is_", "needs_", "should_", "use_", "was_", "with_",
}

var predicateExceptions = map[string]bool{
	"_":     true,
	"ok":    true,
	"value": true,
}
