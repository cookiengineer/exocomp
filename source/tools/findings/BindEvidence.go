package findings

import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"

func (tool *Findings) BindEvidence(id string, arguments utils_chat.Arguments) (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	err0 := readFindingsIndex(tool)

	if err0 == nil {

		finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

		if err1 == nil {

			err2 := bindEvidenceToFinding(tool, &finding, arguments)

			if err2 == nil {

				err3 := writeFinding(tool, finding)

				if err3 == nil {

					evidence := finding.Evidence[len(finding.Evidence)-1]

					return fmt.Sprintf("findings.BindEvidence: Evidence \"%s\" bound to finding \"%s\" (%s).", evidence.ID, finding.ID, evidence.SHA256), nil

				} else {
					return "", fmt.Errorf("findings.BindEvidence: %s", err3.Error())
				}

			} else {
				return "", fmt.Errorf("findings.BindEvidence: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("findings.BindEvidence: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("findings.BindEvidence: %s", err0.Error())
	}

}
