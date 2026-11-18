package findings

import "crypto/rand"
import "encoding/hex"

func newPrefixedID(prefix string) string {

	buffer := make([]byte, 4)
	_, err := rand.Read(buffer)

	if err == nil {
		return prefix + hex.EncodeToString(buffer)
	} else {
		return ""
	}

}
