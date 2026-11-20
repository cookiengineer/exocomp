package types

import "encoding/hex"
import "strings"

func IsAssetID(value string) bool {

	if len(value) == 12 && strings.HasPrefix(value, "ast_") == true {

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
