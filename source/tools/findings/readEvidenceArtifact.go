package findings

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"

func readEvidenceArtifact(tool *Findings, evidence types.FindingEvidence, evidence_id string) ([]byte, string, error) {

	source, err0 := utils_sandbox.ResolvePath(tool.Playground, evidence.Artifact)

	if err0 == nil {

		buffer, err1 := os.ReadFile(source)

		if err1 == nil {
			return buffer, source, nil
		} else {
			return nil, "", fmt.Errorf("findings.Export: Cannot read evidence \"%s\".", evidence_id)
		}

	} else {
		return nil, "", fmt.Errorf("findings.Export: Cannot resolve evidence \"%s\".", evidence_id)
	}

}
