package findings

import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import "fmt"

func readFindingEvidence(tool *Findings, id string, evidence_id string) (types.FindingEvidence, error) {

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 == nil {

		evidence, ok := findEvidence(finding, evidence_id)

		if ok == true && evidence.Artifact != "" {
			return evidence, nil
		} else {
			return types.FindingEvidence{}, fmt.Errorf("findings.Export: Evidence \"%s\" has no artifact.", evidence_id)
		}

	} else {
		return types.FindingEvidence{}, fmt.Errorf("findings.Export: %s", err0.Error())
	}

}
