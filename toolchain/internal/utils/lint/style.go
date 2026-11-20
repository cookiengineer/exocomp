package lint

import "regexp"

var camelCasePattern = regexp.MustCompile(`^[a-z][a-z0-9]*([A-Z][a-z0-9]*)+$`)

var allowedLocalAcronyms = map[string]bool{
	"ANSI":  true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CVE":   true,
	"CVSS":  true,
	"CWE":   true,
	"DNS":   true,
	"EOF":   true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IO":    true,
	"IP":    true,
	"JSON":  true,
	"OS":    true,
	"PTY":   true,
	"SHA":   true,
	"SQL":   true,
	"TLS":   true,
	"TTL":   true,
	"TTY":   true,
	"UI":    true,
	"URI":   true,
	"URL":   true,
	"UUID":  true,
	"YAML":  true,
}
