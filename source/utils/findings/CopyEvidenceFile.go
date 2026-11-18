package findings

import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "io"
import "os"

const maxEvidenceArtifactSize = 4 * 1024 * 1024

func CopyEvidenceFile(playground string, finding_id string, evidence_id string, source string, filename string) (string, string, int64, error) {

	stat, err0 := os.Stat(source)

	if err0 == nil && stat.Mode().IsRegular() == true {

		if stat.Size() <= maxEvidenceArtifactSize {

			absolute, display, err1 := EvidenceArtifactPath(playground, finding_id, evidence_id, filename)

			if err1 == nil {

				source_file, err2 := os.Open(source)

				if err2 == nil {

					defer source_file.Close()

					target_file, err3 := os.Create(absolute)

					if err3 == nil {

						defer target_file.Close()

						hasher := sha256.New()
						written, err4 := io.Copy(io.MultiWriter(target_file, hasher), source_file)

						if err4 == nil {

							err5 := target_file.Sync()

							if err5 == nil {
								return display, hex.EncodeToString(hasher.Sum(nil)), written, nil
							} else {
								return "", "", 0, fmt.Errorf("Cannot persist evidence artifact.")
							}

						} else {
							return "", "", 0, fmt.Errorf("Cannot copy evidence artifact.")
						}

					} else {
						return "", "", 0, fmt.Errorf("Cannot create evidence artifact.")
					}

				} else {
					return "", "", 0, fmt.Errorf("Cannot read evidence artifact.")
				}

			} else {
				return "", "", 0, err1
			}

		} else {
			return "", "", 0, fmt.Errorf("Evidence artifact exceeds the 4 MiB limit.")
		}

	} else if err0 != nil {
		return "", "", 0, fmt.Errorf("Evidence artifact does not exist.")
	} else {
		return "", "", 0, fmt.Errorf("Evidence artifact is not a regular file.")
	}

}
