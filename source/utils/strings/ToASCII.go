package strings

import "encoding/json"
import _ "embed"
import "strings"

//go:embed ToASCII.UTF8.json
var embeddedUTF8Table []byte

var utf8Table map[string][]uint16

func init() {
	json.Unmarshal(embeddedUTF8Table, &utf8Table)
}

func ToASCII(value string) string {

	tmp := strings.TrimSpace(value)

	var filtered string

	for _, character := range tmp {

		var mapped = lookup(uint16(character))

		if mapped != "" {
			filtered += mapped
		}

	}

	return filtered

}
