package findings

import "regexp"

var findingHeaderPattern = regexp.MustCompile(`(?i)(authorization|cookie|set-cookie|x-api-key|api-key)\s*:\s*[^\r\n]+`)

var findingQueryPattern = regexp.MustCompile(`(?i)([?&](?:token|api[_-]?key|apikey|key|secret|password|passwd|session|jwt)=)[^&\s]+`)

var findingUserinfoPattern = regexp.MustCompile(`(?i)(https?://)[^/\s:@]+:[^/\s@]+@`)

var findingBearerPattern = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`)

func RedactSecrets(value string) string {

	value = findingHeaderPattern.ReplaceAllString(value, "$1: [REDACTED]")
	value = findingQueryPattern.ReplaceAllString(value, "$1[REDACTED]")
	value = findingUserinfoPattern.ReplaceAllString(value, "$1[REDACTED]@")
	value = findingBearerPattern.ReplaceAllString(value, "Bearer [REDACTED]")

	return value

}
