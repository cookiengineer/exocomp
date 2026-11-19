package tools

import "exocomp/types"
import utils_findings "exocomp/utils/findings"
import utils_sandbox "exocomp/utils/sandbox"
import "crypto/sha256"
import "encoding/hex"
import "encoding/json"
import "fmt"
import "mime"
import "os"
import "path/filepath"
import "strings"

func (tool *Findings) Export(arguments findingArguments) (string, error) {

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

func exportFindingsData(tool *Findings, export_type string, arguments findingArguments) ([]byte, string, string, error) {

	if export_type == "finding" {
		return exportFindingData(tool, arguments)
	} else if export_type == "evidence" {
		return exportEvidenceData(tool, arguments)
	} else if export_type == "report" {
		return exportReportData(tool)
	} else if export_type == "scan" {
		return exportScanData(tool, arguments)
	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"type\".")
	}

}

func exportFindingData(tool *Findings, arguments findingArguments) ([]byte, string, string, error) {

	id, has_id := arguments.String("id")

	if has_id == true && types.IsFindingID(id) == true {

		finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

		if err0 == nil {

			buffer, err1 := json.MarshalIndent(finding, "", "\t")

			if err1 == nil {
				return buffer, filepath.Join("work", "findings", id+".json"), "application/json", nil
			} else {
				return nil, "", "", fmt.Errorf("findings.Export: Cannot serialize finding \"%s\".", id)
			}

		} else {
			return nil, "", "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
	}

}

func exportEvidenceData(tool *Findings, arguments findingArguments) ([]byte, string, string, error) {

	id, has_id := arguments.String("id")
	evidence_id, has_evidence_id := arguments.String("evidence_id")

	if has_id == true && types.IsFindingID(id) == true {

		if has_evidence_id == true && strings.TrimSpace(evidence_id) != "" {

			evidence, err0 := readFindingEvidence(tool, id, evidence_id)

			if err0 == nil {

				buffer, source, err1 := readEvidenceArtifact(tool, evidence, evidence_id)

				if err1 == nil {
					return buffer, filepath.Join("work", "evidence", id, evidence_id+"-"+filepath.Base(source)), evidence.ContentType, nil
				} else {
					return nil, "", "", err1
				}

			} else {
				return nil, "", "", err0
			}

		} else {
			return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"evidence_id\".")
		}

	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
	}

}

func exportReportData(tool *Findings) ([]byte, string, string, error) {

	report, _, err0 := buildFindingsReport(tool, "", false, "", false)

	if err0 == nil {
		return []byte(report), filepath.Join("work", "findings", "report.md"), "text/markdown", nil
	} else {
		return nil, "", "", fmt.Errorf("findings.Export: %s", err0.Error())
	}

}

func exportScanData(tool *Findings, arguments findingArguments) ([]byte, string, string, error) {

	id, has_id := arguments.String("id")

	if has_id == true && types.IsScanID(id) == true {

		buffer, err0 := utils_findings.ReadScanArtifact(tool.Playground, id)

		if err0 == nil {
			return buffer, filepath.Join("work", "scans", id+"-results.jsonl"), "application/x-ndjson", nil
		} else {
			return nil, "", "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

	} else {
		return nil, "", "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
	}

}

func readFindingEvidence(tool *Findings, id string, evidence_id string) (types.FindingEvidence, error) {

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 == nil {

		evidence, ok := findEvidence(finding, evidence_id)

		if ok == true && evidence.Artifact != "" {
			return evidence, nil
		} else {
			return types.FindingEvidence{}, fmt.Errorf("findings.Export: Evidence \"%s\" has no artifact.", evidence_id)
		}

	} else {
		return types.FindingEvidence{}, fmt.Errorf("findings.Export: %s", err0.Error())
	}

}

func readEvidenceArtifact(tool *Findings, evidence types.FindingEvidence, evidence_id string) ([]byte, string, error) {

	source, err0 := utils_sandbox.ResolvePath(tool.Playground, evidence.Artifact)

	if err0 == nil {

		buffer, err1 := os.ReadFile(source)

		if err1 == nil {
			return buffer, source, nil
		} else {
			return nil, "", fmt.Errorf("findings.Export: Cannot read evidence \"%s\".", evidence_id)
		}

	} else {
		return nil, "", fmt.Errorf("findings.Export: Cannot resolve evidence \"%s\".", evidence_id)
	}

}

func writeFindingsExport(tool *Findings, export_type string, data []byte, default_path string, content_type string, arguments findingArguments) (string, error) {

	path, has_path := arguments.String("path")

	if has_path == true && strings.TrimSpace(path) != "" {
		default_path = path
	}

	sandbox := tool.Sandbox

	if sandbox == "" {
		sandbox = tool.Playground
	}

	if isReservedFilesPath(sandbox, default_path) == true {
		return "", fmt.Errorf("findings.Export: Path is reserved for internal state.")
	}

	resolved, display, err0 := writeFindingsExportFile(sandbox, default_path, data)

	if err0 == nil {

		if content_type == "" {
			content_type = contentTypeForPath(resolved)
		}

		sum := sha256.Sum256(data)

		return fmt.Sprintf(
			"findings.Export: type=%s path=\"%s\" sha256=%s size=%d content_type=%s",
			export_type,
			display,
			hex.EncodeToString(sum[:]),
			len(data),
			content_type,
		), nil

	} else {
		return "", fmt.Errorf("findings.Export: %s", err0.Error())
	}

}

func writeFindingsExportFile(sandbox string, path string, data []byte) (string, string, error) {

	resolved, err0 := utils_sandbox.ResolvePath(sandbox, path)

	if err0 == nil {

		if err1 := os.WriteFile(resolved, data, 0644); err1 == nil {

			display, err2 := utils_sandbox.SanitizePath(sandbox, resolved)

			if err2 == nil {
				return resolved, display, nil
			} else {
				return "", "", fmt.Errorf("Cannot resolve export path.")
			}

		} else {
			return "", "", fmt.Errorf("Cannot write export file.")
		}

	} else {
		return "", "", fmt.Errorf("Cannot resolve export path.")
	}

}

func findEvidence(finding types.Finding, evidence_id string) (types.FindingEvidence, bool) {

	for _, evidence := range finding.Evidence {

		if evidence.ID == evidence_id {
			return evidence, true
		}

	}

	return types.FindingEvidence{}, false

}

func contentTypeForPath(path string) string {

	result := mime.TypeByExtension(filepath.Ext(path))

	if result != "" {
		return result
	}

	extension := strings.ToLower(filepath.Ext(path))

	if extension == ".jsonl" {
		return "application/x-ndjson"
	} else if extension == ".log" || extension == ".txt" {
		return "text/plain"
	}

	return "application/octet-stream"

}
