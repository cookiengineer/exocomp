package fmt

import "strings"

// JoinQuoted joins the given values with ", ", wrapping each value in double quotes.
func JoinQuoted(values []string) string {

	quoted := make([]string, 0)

	for _, value := range values {
		quoted = append(quoted, "\""+value+"\"")
	}

	return strings.Join(quoted, ", ")

}
