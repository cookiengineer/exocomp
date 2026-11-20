package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "path/filepath"
import "strings"

func exportEvidenceData(tool *Findings, arguments utils_chat.Arguments) ([]byte, string, string, error) {

	id, has_id := arguments.String("id")
	evidence_id, has_evidence_id := arguments.String("evidence_id")

	if has_id == true && types.IsFindingID(id) == true {

		if has_evidence_id == true && strings.TrimSpace(evidence_id) != "" {

			evidence, err0 := readFindingEvidence(tool, id, evidence_id)

			if err0 == nil {

				buffer, source, err1 := readEvidenceArtifact(tool, evidence, evidence_id)

				if err1 == nil {
					return buffer, filepath.Join("work", "evidence", id, evidence_id+"-"+filepath.Base(source)), evidence.ContentType, nil
				} else {
					return nil, "", "", err1
				}

			} else {
				return nil, "", "", err0
			}

		} else {
			return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"evidence_id\".")
		}

	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
	}

}
