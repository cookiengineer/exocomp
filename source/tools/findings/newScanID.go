package findings

import "crypto/sha256"
import "encoding/hex"

func newScanID(data []byte) string {

	sum := sha256.Sum256(data)

	return "scn_" + hex.EncodeToString(sum[:])[0:12]

}
