package assets

import "crypto/rand"
import "encoding/hex"

func NewAssetID() string {

	buffer := make([]byte, 4)
	_, err := rand.Read(buffer)

	if err == nil {
		return "ast_" + hex.EncodeToString(buffer)
	} else {
		return ""
	}

}
