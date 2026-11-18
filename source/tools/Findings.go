package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import utils_sandbox "exocomp/utils/sandbox"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"
import "slices"
import "sort"
import "strings"
import "sync"
import "time"

const (
	findingsDefaultLimit     = 50
	findingsMaxLimit         = 200
	findingsMaxTitle         = 200
	findingsMaxIdentifier    = 256
	findingsMaxEndpoint      = 2048
	findingsMaxText          = 16 * 1024
	findingsMaxEvidenceText  = 64 * 1024
	findingsMaxEvidenceCount = 32
	findingsMaxArtifactSize  = 4 * 1024 * 1024
)

type findingArguments map[string]interface{}

type Findings struct {
	Methods    []string
	Playground string
	Sandbox    string
	AgentName  string
	mutex      *sync.RWMutex
	contents   map[string]types.FindingSummary
}

func NewFindings(methods []string, playground string, sandbox string, agent_name string) *Findings {

	tool := &Findings{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		AgentName:  agent_name,
		mutex:      &sync.RWMutex{},
		contents:   make(map[string]types.FindingSummary),
	}

	readFindingsIndex(tool)

	return tool

}

func (tool *Findings) Name() string {
	return "findings"
}

func (tool *Findings) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		parsed := findingArguments(arguments)

		if method == "Add" {

			return tool.Add(parsed)

		} else if method == "Import" {

			return tool.Import(parsed)

		} else if method == "BindEvidence" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.BindEvidence(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Get" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Get(id)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "List" {

			return tool.List(parsed)

		} else if method == "Report" {

			return tool.Report(parsed)

		} else if method == "Retest" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Retest(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Search" {

			return tool.Search(parsed)

		} else if method == "Update" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Update(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Verify" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Verify(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else {
			return "", fmt.Errorf("findings.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("findings.%s: Method not allowed.", method)
	}

}

func (tool *Findings) Add(arguments findingArguments) (string, error) {

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

func (tool *Findings) List(arguments findingArguments) (string, error) {

	severity_filter, has_severity := arguments.String("severity")
	status_filter, has_status := arguments.String("status")
	confidence_filter, has_confidence := arguments.String("confidence")
	endpoint_filter, has_endpoint := arguments.String("endpoint")
	asset_filter, has_asset := arguments.String("asset_id")
	limit := findingsDefaultLimit

	if value, ok := arguments.Int("limit"); ok == true {

		if value > 0 {
			limit = value
		}

	}

	if limit > findingsMaxLimit {
		limit = findingsMaxLimit
	}

	if has_severity == true && types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.List: Invalid parameter \"severity\".")
	}

	if has_status == true && types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.List: Invalid parameter \"status\".")
	}

	if has_confidence == true && types.FindingConfidence(strings.ToLower(strings.TrimSpace(confidence_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.List: Invalid parameter \"confidence\".")
	}

	tool.mutex.Lock()
	err0 := readFindingsIndex(tool)

	if err0 == nil {

		summaries := make([]types.FindingSummary, 0)

		for _, summary := range tool.contents {

			if has_severity == true && string(summary.Severity) != strings.ToLower(strings.TrimSpace(severity_filter)) {
				continue
			} else if has_status == true && string(summary.Status) != strings.ToLower(strings.TrimSpace(status_filter)) {
				continue
			} else if has_confidence == true && string(summary.Confidence) != strings.ToLower(strings.TrimSpace(confidence_filter)) {
				continue
			} else if has_endpoint == true && strings.Contains(strings.ToLower(summary.Endpoint), strings.ToLower(strings.TrimSpace(endpoint_filter))) == false {
				continue
			} else if has_asset == true && slices.Contains(summary.AssetIDs, strings.TrimSpace(asset_filter)) == false {
				continue
			}

			summaries = append(summaries, summary)

		}

		utils_findings.SortSummaries(summaries)

		total := len(summaries)

		if len(summaries) > limit {
			summaries = summaries[0:limit]
		}

		lines := make([]string, 0, len(summaries))

		for _, summary := range summaries {
			lines = append(lines, formatFindingSummary(summary))
		}

		tool.mutex.Unlock()

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("findings.List: %d findings (showing %d).", total, len(lines)))

		for _, line := range lines {
			result = append(result, line)
		}

		return strings.Join(result, "\n"), nil

	} else {
		tool.mutex.Unlock()
		return "", fmt.Errorf("findings.List: %s", err0.Error())
	}

}

func (tool *Findings) Search(arguments findingArguments) (string, error) {

	query, ok := arguments.String("query")

	if ok == true {

		query = strings.ToLower(strings.TrimSpace(query))

		if query == "" {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"query\" is empty.")
		}

		severity_filter, has_severity := arguments.String("severity")
		status_filter, has_status := arguments.String("status")
		limit := findingsDefaultLimit

		if value, ok := arguments.Int("limit"); ok == true {

			if value > 0 {
				limit = value
			}

		}

		if limit > findingsMaxLimit {
			limit = findingsMaxLimit
		}

		if has_severity == true && types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == false {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"severity\".")
		}

		if has_status == true && types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == false {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"status\".")
		}

		tool.mutex.Lock()
		findings, err0 := utils_findings.ReadAll(tool.Playground)

		if err0 == nil {

			summaries := make([]types.FindingSummary, 0)

			for _, finding := range findings {

				if has_severity == true && string(finding.Severity) != strings.ToLower(strings.TrimSpace(severity_filter)) {
					continue
				} else if has_status == true && string(finding.Status) != strings.ToLower(strings.TrimSpace(status_filter)) {
					continue
				} else if utils_findings.MatchesQuery(finding, query) == false {
					continue
				}

				summaries = append(summaries, types.FindingSummaryFromFinding(finding))

			}

			utils_findings.SortSummaries(summaries)

			total := len(summaries)

			if len(summaries) > limit {
				summaries = summaries[0:limit]
			}

			lines := make([]string, 0, len(summaries))

			for _, summary := range summaries {
				lines = append(lines, formatFindingSummary(summary))
			}

			tool.mutex.Unlock()

			result := make([]string, 0)
			result = append(result, fmt.Sprintf("findings.Search: %d matches (showing %d).", total, len(lines)))

			for _, line := range lines {
				result = append(result, line)
			}

			return strings.Join(result, "\n"), nil

		} else {
			tool.mutex.Unlock()
			return "", fmt.Errorf("findings.Search: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("findings.Search: Invalid parameter \"query\" is not a string.")
	}

}

func (tool *Findings) Get(id string) (string, error) {

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 == nil {

		buffer, err1 := json.MarshalIndent(finding, "", "\t")

		if err1 == nil {
			return strings.Join([]string{
				fmt.Sprintf("findings.Get: %s", finding.ID),
				string(buffer),
			}, "\n"), nil
		} else {
			return "", fmt.Errorf("findings.Get: Cannot serialize finding \"%s\".", id)
		}

	} else {
		return "", fmt.Errorf("findings.Get: %s", err0.Error())
	}

}

func (tool *Findings) Update(id string, arguments findingArguments) (string, error) {

	if arguments.Has("status") == true || arguments.Has("confidence") == true || arguments.Has("proof_basis") == true || arguments.Has("evidence") == true {
		return "", fmt.Errorf("findings.Update: Status and evidence can only be changed via findings.Retest, findings.Verify or findings.BindEvidence.")
	}

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

}

func (tool *Findings) BindEvidence(id string, arguments findingArguments) (string, error) {

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

func (tool *Findings) Verify(id string, arguments findingArguments) (string, error) {

	proof_basis_value, ok := arguments.String("proof_basis")

	if ok == true {

		proof_basis := types.FindingProofBasis(strings.ToLower(strings.TrimSpace(proof_basis_value)))

		if proof_basis.IsValid() == false || proof_basis == types.FindingProofBasisStaticUnexecuted {
			return "", fmt.Errorf("findings.Verify: Invalid parameter \"proof_basis\".")
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
		return "", fmt.Errorf("findings.Verify: Invalid parameter \"proof_basis\" is not a string.")
	}

}

func (tool *Findings) Retest(id string, arguments findingArguments) (string, error) {

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

func (tool *Findings) Report(arguments findingArguments) (string, error) {

	severity_filter, has_severity := arguments.String("severity")
	status_filter, has_status := arguments.String("status")

	if has_severity == true && types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.Report: Invalid parameter \"severity\".")
	}

	if has_status == true && types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.Report: Invalid parameter \"status\".")
	}

	tool.mutex.Lock()
	findings, err0 := utils_findings.ReadAll(tool.Playground)
	tool.mutex.Unlock()

	if err0 == nil {

		filtered := make([]types.Finding, 0)

		for _, finding := range findings {

			if has_severity == true && string(finding.Severity) != strings.ToLower(strings.TrimSpace(severity_filter)) {
				continue
			} else if has_status == true && string(finding.Status) != strings.ToLower(strings.TrimSpace(status_filter)) {
				continue
			}

			filtered = append(filtered, finding)

		}

		report := utils_findings.BuildReport(filtered, time.Now().Format("2006-01-02 15:04:05"))
		resolved, err1 := utils_sandbox.ResolvePath(tool.Playground, filepath.Join(".exocomp", "findings", "report.md"))

		if err1 == nil {

			err2 := os.WriteFile(resolved, []byte(report), 0666)

			if err2 == nil {

				display, err3 := utils_sandbox.SanitizePath(tool.Playground, resolved)

				if err3 == nil {
					return fmt.Sprintf("findings.Report: Wrote %s with %d findings.", display, len(filtered)), nil
				} else {
					return "", fmt.Errorf("findings.Report: %s", err3.Error())
				}

			} else {
				return "", fmt.Errorf("findings.Report: Cannot write report.")
			}

		} else {
			return "", fmt.Errorf("findings.Report: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("findings.Report: %s", err0.Error())
	}

}

func (tool *Findings) GetContent(id string) (any, error) {

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 == nil {
		return finding, nil
	} else {
		return nil, fmt.Errorf("findings.GetContent: %s", err0.Error())
	}

}

func (tool *Findings) GetContentIdentifiers() []string {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	readFindingsIndex(tool)

	result := make([]string, 0, len(tool.contents))

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	sort.Strings(result)

	return result

}

func (tool *Findings) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Findings) MarshalJSON() ([]byte, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	err0 := readFindingsIndex(tool)

	if err0 == nil {

		summaries := make([]types.FindingSummary, 0, len(tool.contents))

		for _, summary := range tool.contents {
			summaries = append(summaries, summary)
		}

		utils_findings.SortSummaries(summaries)

		return json.MarshalIndent(summaries, "", "\t")

	} else {
		return nil, fmt.Errorf("findings.MarshalJSON: %s", err0.Error())
	}

}

func (tool *Findings) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range FindingsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}

func (arguments findingArguments) Bool(key string) (bool, bool) {

	value, ok := arguments[key].(bool)

	if ok == true {
		return value, true
	} else {
		return false, false
	}

}

func (arguments findingArguments) Float(key string) (float64, bool) {

	value, ok := arguments[key].(float64)

	if ok == true {
		return value, true
	} else {
		return 0, false
	}

}

func (arguments findingArguments) Has(key string) bool {

	_, ok := arguments[key]

	return ok

}

func (arguments findingArguments) Int(key string) (int, bool) {

	value, ok := arguments[key]

	if ok == true {

		switch typed := value.(type) {
		case float64:
			return int(typed), true
		case int:
			return typed, true
		default:
			return 0, false
		}

	} else {
		return 0, false
	}

}

func (arguments findingArguments) Map(key string) (map[string]string, bool) {

	raw, ok := arguments[key]

	if ok == true {

		switch typed := raw.(type) {
		case map[string]string:
			return typed, true
		case map[string]interface{}:

			result := make(map[string]string, len(typed))

			for name, value := range typed {

				if text, ok := value.(string); ok == true {
					result[name] = text
				}

			}

			return result, true

		default:
			return nil, false
		}

	} else {
		return nil, false
	}

}

func (arguments findingArguments) Object(key string) (findingArguments, bool) {

	raw, ok := arguments[key].(map[string]interface{})

	if ok == true {
		return findingArguments(raw), true
	} else {
		return nil, false
	}

}

func (arguments findingArguments) String(key string) (string, bool) {

	value, ok := arguments[key].(string)

	if ok == true {
		return value, true
	} else {
		return "", false
	}

}

func (arguments findingArguments) Strings(key string) ([]string, bool) {

	raw, ok := arguments[key]

	if ok == true {

		switch typed := raw.(type) {
		case []string:
			return typed, true
		case []interface{}:

			result := make([]string, 0, len(typed))

			for _, value := range typed {

				if text, ok := value.(string); ok == true {
					result = append(result, text)
				}

			}

			return result, true

		default:
			return nil, false
		}

	} else {
		return nil, false
	}

}
