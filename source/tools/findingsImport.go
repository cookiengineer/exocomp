package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import "crypto/sha256"
import "encoding/hex"
import "encoding/json"
import "fmt"
import "strings"

type nucleiImportResult struct {
	Finding  types.Finding
	Request  string
	Response string
}

type nucleiImportClassification struct {
	CVEID     []string `json:"cve-id"`
	CWEID     []string `json:"cwe-id"`
	CVSSScore float64  `json:"cvss-score"`
}

type nucleiImportInfo struct {
	Name           string                     `json:"name"`
	Severity       string                     `json:"severity"`
	Description    string                     `json:"description"`
	Reference      []string                   `json:"reference"`
	Classification nucleiImportClassification `json:"classification"`
}

type nucleiImportLine struct {
	TemplateID       string           `json:"template-id"`
	TemplatePath     string           `json:"template-path"`
	Info             nucleiImportInfo `json:"info"`
	Type             string           `json:"type"`
	Host             string           `json:"host"`
	URL              string           `json:"url"`
	MatchedAt        string           `json:"matched-at"`
	Request          string           `json:"request"`
	Response         string           `json:"response"`
	MatcherName      string           `json:"matcher-name"`
	ExtractedResults []string         `json:"extracted-results"`
}

func (tool *Findings) Import(arguments findingArguments) (string, error) {

	path, ok := arguments.String("path")

	if ok == false || strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("findings.Import: Invalid parameter \"path\" is not a non-empty string.")
	}

	format := "nuclei"

	if value, ok := arguments.String("format"); ok == true {
		format = strings.ToLower(strings.TrimSpace(value))
	}

	if format != "nuclei" {
		return "", fmt.Errorf("findings.Import: Unsupported format \"%s\".", format)
	}

	command, _ := arguments.String("command")

	sandbox := tool.Sandbox

	if sandbox == "" {
		sandbox = tool.Playground
	}

	data, err0 := readScanArtifactFile(sandbox, path, "findings.Import")

	if err0 == nil {

		scan_id, has_scan_id := arguments.String("scan_id")

		if has_scan_id == false || strings.TrimSpace(scan_id) == "" {
			scan_id = newScanID(data)
		} else {
			scan_id = strings.TrimSpace(scan_id)
		}

		if types.IsScanID(scan_id) == true {

			if err_store := utils_findings.WriteScanArtifact(tool.Playground, scan_id, data); err_store == nil {

				parsed, invalid, err1 := parseNucleiResults(data, scan_id)

				if err1 == nil {

					imported, skipped, invalid, err2 := importFindingsResults(tool, parsed, invalid, command)

					if err2 == nil {
						return fmt.Sprintf(
							"findings.Import: Imported %d candidates, skipped %d duplicates, %d invalid results (scan %s).",
							imported,
							skipped,
							invalid,
							scan_id,
						), nil
					} else {
						return "", fmt.Errorf("findings.Import: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("findings.Import: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("findings.Import: %s", err_store.Error())
			}

		} else {
			return "", fmt.Errorf("findings.Import: Invalid parameter \"scan_id\".")
		}

	} else {
		return "", err0
	}

}

func importFindingsResults(tool *Findings, parsed []nucleiImportResult, invalid int, command string) (int, int, int, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err0 := readFindingsIndex(tool); err0 == nil {

		existing, err1 := utils_findings.ReadAll(tool.Playground)

		if err1 == nil {

			seen := make(map[string]bool, len(existing))

			for _, finding := range existing {
				seen[importDedupKey(finding)] = true
			}

			imported := 0
			skipped := 0

			for _, result := range parsed {

				if seen[importDedupKey(result.Finding)] == true {

					skipped++
					continue

				}

				if result.Request != "" {

					evidence := findingArguments{
						"kind":    string(types.FindingEvidenceKindHTTPRequest),
						"title":   "Nuclei request",
						"source":  "nuclei",
						"content": result.Request,
					}

					if command != "" {
						evidence["command"] = command
					}

					if err2 := bindEvidenceToFinding(tool, &result.Finding, evidence); err2 != nil {
						invalid++
						continue
					}

				}

				if result.Response != "" {

					evidence := findingArguments{
						"kind":    string(types.FindingEvidenceKindHTTPResponse),
						"title":   "Nuclei response",
						"source":  "nuclei",
						"content": result.Response,
					}

					if command != "" {
						evidence["command"] = command
					}

					if err3 := bindEvidenceToFinding(tool, &result.Finding, evidence); err3 != nil {
						invalid++
						continue
					}

				}

				if result.Finding.Validate() == true {

					if err4 := writeFinding(tool, result.Finding); err4 == nil {
						seen[importDedupKey(result.Finding)] = true
						imported++
					} else {
						return 0, 0, 0, err4
					}

				} else {
					invalid++
				}

			}

			return imported, skipped, invalid, nil

		} else {
			return 0, 0, 0, err1
		}

	} else {
		return 0, 0, 0, err0
	}

}

func parseNucleiResults(data []byte, scan_id string) ([]nucleiImportResult, int, error) {

	trimmed := strings.TrimSpace(string(data))

	if trimmed == "" {
		return make([]nucleiImportResult, 0), 0, nil
	}

	lines := make([]json.RawMessage, 0)

	if strings.HasPrefix(trimmed, "[") == true {

		raw := make([]json.RawMessage, 0)

		if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
			return nil, 0, fmt.Errorf("Cannot parse nuclei results.")
		}

		lines = append(lines, raw...)

	} else {

		for _, line := range strings.Split(trimmed, "\n") {

			line = strings.TrimSpace(line)

			if line != "" {
				lines = append(lines, json.RawMessage(line))
			}

		}

	}

	result := make([]nucleiImportResult, 0, len(lines))
	invalid := 0

	for _, line := range lines {

		if strings.HasPrefix(string(line), "#") == true {
			continue
		}

		parsed := nucleiImportLine{}
		err := json.Unmarshal(line, &parsed)

		if err != nil {
			invalid++
			continue
		}

		if parsed.TemplateID == "" && parsed.Info.Name == "" && parsed.MatchedAt == "" && parsed.Host == "" {
			invalid++
			continue
		}

		result = append(result, nucleiImportLineToResult(parsed, scan_id))

	}

	return result, invalid, nil

}

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

func nucleiReproduction(line nucleiImportLine, endpoint string) string {

	parts := make([]string, 0)

	if line.TemplateID != "" {
		parts = append(parts, fmt.Sprintf("Nuclei template \"%s\" matched.", line.TemplateID))
	} else {
		parts = append(parts, "Nuclei template matched.")
	}

	if endpoint != "" {
		parts = append(parts, fmt.Sprintf("Endpoint: %s", endpoint))
	}

	if line.MatcherName != "" {
		parts = append(parts, fmt.Sprintf("Matcher: %s", line.MatcherName))
	}

	if len(line.ExtractedResults) > 0 {
		parts = append(parts, fmt.Sprintf("Extracted: %s", strings.Join(line.ExtractedResults, ", ")))
	}

	return utils_findings.LimitText(utils_findings.RedactSecrets(strings.Join(parts, " ")), findingsMaxText)

}

func newScanID(data []byte) string {

	sum := sha256.Sum256(data)

	return "scn_" + hex.EncodeToString(sum[:])[0:12]

}

func importDedupKey(finding types.Finding) string {
	return strings.ToLower(strings.Join([]string{finding.Title, finding.Endpoint, finding.Parameter}, "\x00"))
}
