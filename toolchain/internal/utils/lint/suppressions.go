package lint

import "regexp"

var nolintPattern = regexp.MustCompile(`^//\s*nolint:\s*([a-z0-9/,._-]+)`)
