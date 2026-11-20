package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Findings) Verify(id string, arguments utils_chat.Arguments) (string, error) {

	proof_basis_value, ok := arguments.String("proof_basis")

	if ok == true {

		proof_basis := types.FindingProofBasis(strings.ToLower(strings.TrimSpace(proof_basis_value)))

		if proof_basis.IsValid() == true && proof_basis != types.FindingProofBasisStaticUnexecuted {

			tool.mutex.Lock()
			defer tool.mutex.Unlock()

			err0 := readFindingsIndex(tool)

			if err0 == nil {

				finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

				if err1 == nil {

					if evidence_arguments, ok := arguments.Object("evidence"); ok == true {

						err2 := bindEvidenceToFinding(tool, &finding, evidence_arguments)

						if err2 != nil {
							return "", fmt.Errorf("findings.Verify: %s", err2.Error())
						}

					}

					require_reproduction, has_require_reproduction := arguments.Bool("require_reproduction")

					if has_require_reproduction == true && require_reproduction == true {

						has_reproduction := false
						has_test_run := false
						has_environment := finding.Preconditions != ""

						for _, evidence := range finding.Evidence {

							if evidence.Kind.IsReproduction() == true {
								has_reproduction = true
							}

							if evidence.Kind == types.FindingEvidenceKindTestRun && evidence.ExitCode != nil && *evidence.ExitCode == 0 {
								has_test_run = true
							}

							if len(evidence.Environment) > 0 {
								has_environment = true
							}

						}

						if has_reproduction == false || has_test_run == false || has_environment == false {
							return "", fmt.Errorf("findings.Verify: Reproduction requires a unit_test or poc artifact, a successful test_run, and environment or preconditions.")
						}

					} else if len(finding.Evidence) == 0 {
						return "", fmt.Errorf("findings.Verify: Finding \"%s\" has no evidence.", finding.ID)
					}

					finding.Status = types.FindingStatusVerified
					finding.Confidence = types.FindingConfidenceObserved
					finding.ProofBasis = proof_basis
					finding.VerifiedAt = schemas.NewDatetime()
					finding.UpdatedAt = schemas.NewDatetime()

					if tool.AgentName != "" {
						finding.VerifiedBy = utils_findings.LimitText(tool.AgentName, findingsMaxTitle)
					}

					if finding.Validate() == true {

						err3 := writeFinding(tool, finding)

						if err3 == nil {
							return fmt.Sprintf("findings.Verify: Finding \"%s\" verified via %s.", finding.ID, finding.ProofBasis), nil
						} else {
							return "", fmt.Errorf("findings.Verify: %s", err3.Error())
						}

					} else {
						return "", fmt.Errorf("findings.Verify: Invalid finding values.")
					}

				} else {
					return "", fmt.Errorf("findings.Verify: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("findings.Verify: %s", err0.Error())
			}

		} else {
			return "", fmt.Errorf("findings.Verify: Invalid parameter \"proof_basis\".")
		}
	} else {
		return "", fmt.Errorf("findings.Verify: Invalid parameter \"proof_basis\" is not a string.")
	}

}
