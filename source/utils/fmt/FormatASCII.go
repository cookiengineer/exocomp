package fmt

import "encoding/json"
import _ "embed"
import "strings"

//go:embed FormatASCII.json
var embeddedUTF8Table []byte

var utf8Table map[string][]uint16

func init() {
	json.Unmarshal(embeddedUTF8Table, &utf8Table)
}

func FormatASCII(value string) string {

	tmp := strings.TrimSpace(value)
	filtered := ""

	for _, character := range tmp {

		mapped := lookup(uint16(character))

		if mapped != "" {
			filtered += mapped
		}

	}

	return filtered

}
