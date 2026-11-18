package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import utils_sandbox "exocomp/utils/sandbox"
import "crypto/sha256"
import "encoding/hex"
import "encoding/json"
import "fmt"
import "os"
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

	resolved, err0 := utils_sandbox.ResolvePath(sandbox, path)

	if err0 != nil {
		return "", fmt.Errorf("findings.Import: Cannot resolve artifact \"%s\".", path)
	}

	stat, err1 := os.Stat(resolved)

	if err1 != nil || stat.Mode().IsRegular() == false {
		return "", fmt.Errorf("findings.Import: Artifact \"%s\" does not exist.", path)
	}

	if stat.Size() > maxScanArtifactSize {
		return "", fmt.Errorf("findings.Import: Artifact exceeds the 16 MiB limit.")
	}

	data, err2 := os.ReadFile(resolved)

	if err2 != nil {
		return "", fmt.Errorf("findings.Import: Cannot read artifact \"%s\".", path)
	}

	scan_id, has_scan_id := arguments.String("scan_id")

	if has_scan_id == false || strings.TrimSpace(scan_id) == "" {
		scan_id = newScanID(data)
	} else {
		scan_id = strings.TrimSpace(scan_id)
	}

	if types.IsScanID(scan_id) == false {
		return "", fmt.Errorf("findings.Import: Invalid parameter \"scan_id\".")
	}

	if err_store := utils_findings.WriteScanArtifact(tool.Playground, scan_id, data); err_store != nil {
		return "", fmt.Errorf("findings.Import: %s", err_store.Error())
	}

	parsed, invalid, err3 := parseNucleiResults(data, scan_id)

	if err3 != nil {
		return "", fmt.Errorf("findings.Import: %s", err3.Error())
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	err4 := readFindingsIndex(tool)

	if err4 != nil {
		return "", fmt.Errorf("findings.Import: %s", err4.Error())
	}

	existing, err5 := utils_findings.ReadAll(tool.Playground)

	if err5 != nil {
		return "", fmt.Errorf("findings.Import: %s", err5.Error())
	}

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

			if err6 := bindEvidenceToFinding(tool, &result.Finding, evidence); err6 != nil {
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

			if err7 := bindEvidenceToFinding(tool, &result.Finding, evidence); err7 != nil {
				invalid++
				continue
			}

		}

		if result.Finding.Validate() == false {
			invalid++
			continue
		}

		if err8 := writeFinding(tool, result.Finding); err8 != nil {
			return "", fmt.Errorf("findings.Import: %s", err8.Error())
		}

		seen[importDedupKey(result.Finding)] = true
		imported++

	}

	return fmt.Sprintf(
		"findings.Import: Imported %d candidates, skipped %d duplicates, %d invalid results (scan %s).",
		imported,
		skipped,
		invalid,
		scan_id,
	), nil

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
