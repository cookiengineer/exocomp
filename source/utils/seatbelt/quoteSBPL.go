package seatbelt

import "fmt"
import "strings"

func quoteSBPL(value string) string {

	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)

	return fmt.Sprintf("\"%s\"", escaped)

}
