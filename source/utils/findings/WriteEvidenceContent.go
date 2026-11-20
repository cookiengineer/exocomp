package findings

import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "os"

func WriteEvidenceContent(playground string, finding_id string, evidence_id string, filename string, content string) (string, string, int64, error) {

	absolute, display, err0 := EvidenceArtifactPath(playground, finding_id, evidence_id, filename)

	if err0 == nil {

		data := []byte(content)
		err1 := os.WriteFile(absolute, data, 0644)

		if err1 == nil {

			sum := sha256.Sum256(data)

			return display, hex.EncodeToString(sum[:]), int64(len(data)), nil

		} else {
			return "", "", 0, fmt.Errorf("findings.WriteEvidenceContent: Cannot write evidence artifact.")
		}

	} else {
		return "", "", 0, err0
	}

}
