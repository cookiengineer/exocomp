package types

import "encoding/hex"
import "strings"

func IsScanID(value string) bool {

	if len(value) == 16 && strings.HasPrefix(value, "scn_") == true {

		_, err := hex.DecodeString(value[4:])

		if err == nil {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}
