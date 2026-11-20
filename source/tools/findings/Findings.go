package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
import "strings"
import "sync"

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

func (tool *Findings) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		parsed := utils_chat.Arguments(arguments)

		if method == "Add" {

			return tool.Add(parsed)

		} else if method == "Export" {

			return tool.Export(parsed)

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

func (tool *Findings) Export(arguments utils_chat.Arguments) (string, error) {

	type_value, ok := arguments.String("type")

	if ok == true {

		export_type := strings.ToLower(strings.TrimSpace(type_value))
		data, default_path, content_type, err0 := exportFindingsData(tool, export_type, arguments)

		if err0 == nil {
			return writeFindingsExport(tool, export_type, data, default_path, content_type, arguments)
		} else {
			return "", err0
		}

	} else {
		return "", fmt.Errorf("findings.Export: Invalid parameter \"type\" is not a string.")
	}

}

func (tool *Findings) Get(id string) (string, error) {

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 == nil {

		view := newFindingView(finding)
		buffer, err1 := json.MarshalIndent(view, "", "\t")

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

func (tool *Findings) Import(arguments utils_chat.Arguments) (string, error) {

	path, ok := arguments.String("path")

	if ok == true && strings.TrimSpace(path) != "" {

		format := "nuclei"

		if value, ok := arguments.String("format"); ok == true {
			format = strings.ToLower(strings.TrimSpace(value))
		}

		if format == "nuclei" {

			command, _ := arguments.String("command")

			sandbox := tool.Sandbox

			if sandbox == "" {
				sandbox = tool.Playground
			}

			data, err0 := utils_findings.ReadScanArtifactFile(sandbox, path, "findings.Import")

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

		} else {
			return "", fmt.Errorf("findings.Import: Unsupported format \"%s\".", format)
		}

	} else {
		return "", fmt.Errorf("findings.Import: Invalid parameter \"path\" is not a non-empty string.")
	}

}

func (tool *Findings) List(arguments utils_chat.Arguments) (string, error) {

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

	if has_severity == false || types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == true {

		if has_status == false || types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == true {

			if has_confidence == false || types.FindingConfidence(strings.ToLower(strings.TrimSpace(confidence_filter))).IsValid() == true {

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

			} else {
				return "", fmt.Errorf("findings.List: Invalid parameter \"confidence\".")
			}
		} else {
			return "", fmt.Errorf("findings.List: Invalid parameter \"status\".")
		}
	} else {
		return "", fmt.Errorf("findings.List: Invalid parameter \"severity\".")
	}

}

func (tool *Findings) Report(arguments utils_chat.Arguments) (string, error) {

	severity_filter, has_severity := arguments.String("severity")
	status_filter, has_status := arguments.String("status")

	if has_severity == false || types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == true {

		if has_status == false || types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == true {

			report, count, err0 := buildFindingsReport(tool, severity_filter, has_severity, status_filter, has_status)

			if err0 == nil {
				return fmt.Sprintf("findings.Report: %d findings.\n%s", count, report), nil
			} else {
				return "", fmt.Errorf("findings.Report: %s", err0.Error())
			}

		} else {
			return "", fmt.Errorf("findings.Report: Invalid parameter \"status\".")
		}
	} else {
		return "", fmt.Errorf("findings.Report: Invalid parameter \"severity\".")
	}

}

func (tool *Findings) Retest(id string, arguments utils_chat.Arguments) (string, error) {

	status_value, ok := arguments.String("status")

	if ok == true {

		status := types.FindingStatus(strings.ToLower(strings.TrimSpace(status_value)))

		if status.IsValid() == true && status != types.FindingStatusCandidate {

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
			return "", fmt.Errorf("findings.Retest: Invalid parameter \"status\".")
		}

	} else {
		return "", fmt.Errorf("findings.Retest: Invalid parameter \"status\" is not a string.")
	}

}

func (tool *Findings) Search(arguments utils_chat.Arguments) (string, error) {

	query, ok := arguments.String("query")

	if ok == true {

		query = strings.ToLower(strings.TrimSpace(query))

		if query != "" {

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

			if has_severity == false || types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == true {

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
				return "", fmt.Errorf("findings.Search: Invalid parameter \"severity\".")
			}

		} else {
			return "", fmt.Errorf("findings.Search: Invalid parameter \"query\" is empty.")
		}

	} else {
		return "", fmt.Errorf("findings.Search: Invalid parameter \"query\" is not a string.")
	}

}

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
