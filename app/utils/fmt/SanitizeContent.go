package fmt

import "strings"

var sanitizeReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&quot;",
	"'", "&#39;",
)

func SanitizeContent(raw string) string {
	return sanitizeReplacer.Replace(raw)
}
