package lint

import "regexp"

var prefixedErrorPattern = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9_]*\.)+(%s|%v|[A-Z][A-Za-z0-9]*):`)

var parameterPrefixPattern = regexp.MustCompile(`^%(s|v)([.:]|$)`)
