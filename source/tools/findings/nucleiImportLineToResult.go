package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import "strings"

func nucleiImportLineToResult(line nucleiImportLine, scan_id string) nucleiImportResult {

	now := schemas.NewDatetime()

	title := strings.TrimSpace(line.Info.Name)

	if title == "" {
		title = strings.TrimSpace(line.TemplateID)
	}

	if title == "" {
		title = "Nuclei finding"
	}

	endpoint := strings.TrimSpace(line.MatchedAt)

	if endpoint == "" {
		endpoint = strings.TrimSpace(line.URL)
	}

	if endpoint == "" {
		endpoint = strings.TrimSpace(line.Host)
	}

	severity := types.FindingSeverity(strings.ToLower(strings.TrimSpace(line.Info.Severity)))

	if severity.IsValid() == false {
		severity = types.FindingSeverityInfo
	}

	finding := types.Finding{
		ID:           utils_findings.NewFindingID(),
		Title:        utils_findings.LimitText(utils_findings.RedactSecrets(title), findingsMaxTitle),
		Severity:     severity,
		AssetIDs:     make([]string, 0),
		Endpoint:     utils_findings.LimitText(utils_findings.RedactSecrets(endpoint), findingsMaxEndpoint),
		Status:       types.FindingStatusCandidate,
		Confidence:   types.FindingConfidenceUnverified,
		ProofBasis:   types.FindingProofBasisStaticUnexecuted,
		Evidence:     make([]types.FindingEvidence, 0),
		SourceScanID: utils_findings.LimitText(scan_id, findingsMaxIdentifier),
		Reproduction: nucleiReproduction(line, endpoint),
		Impact:       utils_findings.LimitText(utils_findings.RedactSecrets(line.Info.Description), findingsMaxText),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if len(line.Info.Classification.CWEID) > 0 {
		finding.CWE = utils_findings.LimitText(line.Info.Classification.CWEID[0], findingsMaxIdentifier)
	}

	if line.Info.Classification.CVSSScore > 0 && line.Info.Classification.CVSSScore <= 10 {
		finding.CVSS = line.Info.Classification.CVSSScore
	}

	return nucleiImportResult{
		Finding:  finding,
		Request:  line.Request,
		Response: line.Response,
	}

}
