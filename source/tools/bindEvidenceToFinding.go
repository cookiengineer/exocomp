package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "path/filepath"

func bindEvidenceToFinding(tool *Findings, finding *types.Finding, arguments findingArguments) error {

	if len(finding.Evidence) < findingsMaxEvidenceCount {

		evidence, provided, err0 := parseFindingEvidence(arguments)

		if err0 == nil && provided == true {

			evidence.ID = utils_findings.NewEvidenceID()
			evidence.CreatedAt = schemas.NewDatetime()

			if evidence.Content != "" {

				evidence.Content = utils_findings.LimitText(utils_findings.RedactSecrets(evidence.Content), findingsMaxEvidenceText)

				filename := "evidence.txt"

				if evidence.ContentType == "" {
					evidence.ContentType = "text/plain"
				}

				artifact, sha256_value, size, err1 := utils_findings.WriteEvidenceContent(tool.Playground, finding.ID, evidence.ID, filename, evidence.Content)

				if err1 == nil {
					evidence.Artifact = artifact
					evidence.SHA256 = sha256_value
					evidence.Size = size
				} else {
					return fmt.Errorf("findings.BindEvidence: %s", err1.Error())
				}

			} else {

				resolved, err1 := utils_sandbox.ResolvePath(tool.Sandbox, evidence.Artifact)

				if err1 == nil {

					artifact, sha256_value, size, err2 := utils_findings.CopyEvidenceFile(tool.Playground, finding.ID, evidence.ID, resolved, filepath.Base(resolved))

					if err2 == nil {
						evidence.Artifact = artifact
						evidence.SHA256 = sha256_value
						evidence.Size = size
					} else {
						return fmt.Errorf("findings.BindEvidence: %s", err2.Error())
					}

				} else {
					return fmt.Errorf("findings.BindEvidence: %s", err1.Error())
				}

			}

			finding.Evidence = append(finding.Evidence, evidence)
			finding.UpdatedAt = schemas.NewDatetime()

			return nil

		} else if err0 != nil {
			return err0
		} else {
			return fmt.Errorf("findings.BindEvidence: Evidence is required.")
		}

	} else {
		return fmt.Errorf("findings.BindEvidence: Finding \"%s\" already has %d evidence entries.", finding.ID, findingsMaxEvidenceCount)
	}

}
