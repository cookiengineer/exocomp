package findings

import "exocomp/types"

func findEvidence(finding types.Finding, evidence_id string) (types.FindingEvidence, bool) {

	for _, evidence := range finding.Evidence {

		if evidence.ID == evidence_id {
			return evidence, true
		}

	}

	return types.FindingEvidence{}, false

}
