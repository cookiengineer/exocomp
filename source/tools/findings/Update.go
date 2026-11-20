package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Findings) Update(id string, arguments utils_chat.Arguments) (string, error) {

	if arguments.Has("status") == false && arguments.Has("confidence") == false && arguments.Has("proof_basis") == false && arguments.Has("evidence") == false {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		err0 := readFindingsIndex(tool)

		if err0 == nil {

			finding, err1 := utils_findings.ReadFinding(tool.Playground, id)

			if err1 == nil {

				if title, ok := arguments.String("title"); ok == true {
					finding.Title = utils_findings.LimitText(utils_findings.RedactSecrets(title), findingsMaxTitle)
				}

				if severity, ok := arguments.String("severity"); ok == true {
					finding.Severity = types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity)))
				}

				if cvss, ok := arguments.Float("cvss"); ok == true {
					finding.CVSS = cvss
				}

				if cwe, ok := arguments.String("cwe"); ok == true {
					finding.CWE = utils_findings.LimitText(utils_findings.RedactSecrets(cwe), findingsMaxIdentifier)
				}

				if owasp, ok := arguments.String("owasp"); ok == true {
					finding.OWASP = utils_findings.LimitText(utils_findings.RedactSecrets(owasp), findingsMaxIdentifier)
				}

				if asset_ids, ok := arguments.Strings("asset_ids"); ok == true {
					finding.AssetIDs = utils_findings.LimitIdentifiers(asset_ids, findingsMaxIdentifier)
				}

				if endpoint, ok := arguments.String("endpoint"); ok == true {
					finding.Endpoint = utils_findings.LimitText(utils_findings.RedactSecrets(endpoint), findingsMaxEndpoint)
				}

				if parameter, ok := arguments.String("parameter"); ok == true {
					finding.Parameter = utils_findings.LimitText(utils_findings.RedactSecrets(parameter), findingsMaxIdentifier)
				}

				if preconditions, ok := arguments.String("preconditions"); ok == true {
					finding.Preconditions = utils_findings.LimitText(utils_findings.RedactSecrets(preconditions), findingsMaxText)
				}

				if reproduction, ok := arguments.String("reproduction"); ok == true {
					finding.Reproduction = utils_findings.LimitText(utils_findings.RedactSecrets(reproduction), findingsMaxText)
				}

				if impact, ok := arguments.String("impact"); ok == true {
					finding.Impact = utils_findings.LimitText(utils_findings.RedactSecrets(impact), findingsMaxText)
				}

				if remediation, ok := arguments.String("remediation"); ok == true {
					finding.Remediation = utils_findings.LimitText(utils_findings.RedactSecrets(remediation), findingsMaxText)
				}

				if counterevidence, ok := arguments.String("counterevidence"); ok == true {
					finding.Counterevidence = utils_findings.LimitText(utils_findings.RedactSecrets(counterevidence), findingsMaxText)
				}

				finding.UpdatedAt = schemas.NewDatetime()

				if finding.Validate() == true {

					err2 := writeFinding(tool, finding)

					if err2 == nil {
						return fmt.Sprintf("findings.Update: Finding \"%s\" updated.", finding.ID), nil
					} else {
						return "", fmt.Errorf("findings.Update: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("findings.Update: Invalid finding values.")
				}

			} else {
				return "", fmt.Errorf("findings.Update: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("findings.Update: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("findings.Update: Status and evidence can only be changed via findings.Retest, findings.Verify or findings.BindEvidence.")
	}
}
