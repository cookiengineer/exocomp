package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Findings) Add(arguments utils_chat.Arguments) (string, error) {

	title, ok1 := arguments.String("title")
	severity_value, ok2 := arguments.String("severity")

	if ok1 == true && ok2 == true {

		now := schemas.NewDatetime()

		finding := types.Finding{
			ID:         utils_findings.NewFindingID(),
			Title:      utils_findings.LimitText(utils_findings.RedactSecrets(title), findingsMaxTitle),
			Severity:   types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_value))),
			Status:     types.FindingStatusCandidate,
			Confidence: types.FindingConfidenceUnverified,
			ProofBasis: types.FindingProofBasisStaticUnexecuted,
			AssetIDs:   make([]string, 0),
			Evidence:   make([]types.FindingEvidence, 0),
			CreatedAt:  now,
			UpdatedAt:  now,
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

		if source_scan_id, ok := arguments.String("source_scan_id"); ok == true {
			finding.SourceScanID = utils_findings.LimitText(source_scan_id, findingsMaxIdentifier)
		}

		if finding.Validate() == true {

			tool.mutex.Lock()
			err0 := readFindingsIndex(tool)

			if err0 == nil {
				err0 = writeFinding(tool, finding)
			}

			tool.mutex.Unlock()

			if err0 == nil {
				return fmt.Sprintf("findings.Add: Finding \"%s\" created.", finding.ID), nil
			} else {
				return "", fmt.Errorf("findings.Add: %s", err0.Error())
			}

		} else {
			return "", fmt.Errorf("findings.Add: Invalid finding.")
		}

	} else if ok1 == true && ok2 == false {
		return "", fmt.Errorf("findings.Add: Invalid parameter \"severity\" is not a string.")
	} else if ok1 == false && ok2 == true {
		return "", fmt.Errorf("findings.Add: Invalid parameter \"title\" is not a string.")
	} else {
		return "", fmt.Errorf("findings.Add: Invalid parameters.")
	}

}
