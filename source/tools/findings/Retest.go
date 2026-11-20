package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Findings) Retest(id string, arguments utils_chat.Arguments) (string, error) {

	status_value, ok := arguments.String("status")

	if ok == true {

		status := types.FindingStatus(strings.ToLower(strings.TrimSpace(status_value)))

		if status.IsValid() == false || status == types.FindingStatusCandidate {
			return "", fmt.Errorf("findings.Retest: Invalid parameter \"status\".")
		}

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		err0 := readFindingsIndex(tool)

		if err0 == nil {

			finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

			if err1 == nil {

				if evidence_arguments, ok := arguments.Object("evidence"); ok == true {

					err2 := bindEvidenceToFinding(tool, &finding, evidence_arguments)

					if err2 != nil {
						return "", fmt.Errorf("findings.Retest: %s", err2.Error())
					}

				}

				if counterevidence, ok := arguments.String("counterevidence"); ok == true {
					finding.Counterevidence = utils_findings.LimitText(utils_findings.RedactSecrets(counterevidence), findingsMaxText)
				}

				if proof_basis_value, ok := arguments.String("proof_basis"); ok == true {

					proof_basis := types.FindingProofBasis(strings.ToLower(strings.TrimSpace(proof_basis_value)))

					if proof_basis.IsValid() == false {
						return "", fmt.Errorf("findings.Retest: Invalid parameter \"proof_basis\".")
					}

					finding.ProofBasis = proof_basis

				}

				if (status == types.FindingStatusFixed || status == types.FindingStatusFalsePositive) && len(finding.Evidence) == 0 && finding.Counterevidence == "" {
					return "", fmt.Errorf("findings.Retest: Status \"%s\" requires evidence or counterevidence.", status)
				}

				if status == types.FindingStatusVerified {

					if finding.ProofBasis != types.FindingProofBasisExecutedLive && finding.ProofBasis != types.FindingProofBasisExecutedLab {
						return "", fmt.Errorf("findings.Retest: Verified findings require executed_live or executed_lab proof.")
					} else if len(finding.Evidence) == 0 {
						return "", fmt.Errorf("findings.Retest: Verified findings require evidence.")
					}

					finding.Confidence = types.FindingConfidenceObserved
					finding.VerifiedAt = schemas.NewDatetime()

					if tool.AgentName != "" {
						finding.VerifiedBy = utils_findings.LimitText(tool.AgentName, findingsMaxTitle)
					}

				}

				finding.Status = status
				finding.UpdatedAt = schemas.NewDatetime()

				if finding.Validate() == true {

					err2 := writeFinding(tool, finding)

					if err2 == nil {
						return fmt.Sprintf("findings.Retest: Finding \"%s\" marked as %s.", finding.ID, finding.Status), nil
					} else {
						return "", fmt.Errorf("findings.Retest: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("findings.Retest: Invalid finding values.")
				}

			} else {
				return "", fmt.Errorf("findings.Retest: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("findings.Retest: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("findings.Retest: Invalid parameter \"status\" is not a string.")
	}

}
