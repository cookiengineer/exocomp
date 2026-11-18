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

	if ok == false {
		return "", fmt.Errorf("findings.Export: Invalid parameter \"type\" is not a string.")
	}

	export_type := strings.ToLower(strings.TrimSpace(type_value))
	id, has_id := arguments.String("id")
	evidence_id, has_evidence_id := arguments.String("evidence_id")
	path, has_path := arguments.String("path")

	data := make([]byte, 0)
	default_path := ""
	content_type := ""

	if export_type == "finding" {

		if has_id == false || types.IsFindingID(id) == false {
			return "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
		}

		finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

		if err0 != nil {
			return "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

		buffer, err1 := json.MarshalIndent(finding, "", "\t")

		if err1 != nil {
			return "", fmt.Errorf("findings.Export: Cannot serialize finding \"%s\".", id)
		}

		data = buffer
		default_path = filepath.Join("work", "findings", id+".json")
		content_type = "application/json"

	} else if export_type == "evidence" {

		if has_id == false || types.IsFindingID(id) == false {
			return "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
		} else if has_evidence_id == false || strings.TrimSpace(evidence_id) == "" {
			return "", fmt.Errorf("findings.Export: Invalid parameter \"evidence_id\".")
		}

		finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

		if err0 != nil {
			return "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

		evidence, ok := findEvidence(finding, evidence_id)

		if ok == false || evidence.Artifact == "" {
			return "", fmt.Errorf("findings.Export: Evidence \"%s\" has no artifact.", evidence_id)
		}

		source, err1 := utils_sandbox.ResolvePath(tool.Playground, evidence.Artifact)

		if err1 != nil {
			return "", fmt.Errorf("findings.Export: Cannot resolve evidence \"%s\".", evidence_id)
		}

		buffer, err2 := os.ReadFile(source)

		if err2 != nil {
			return "", fmt.Errorf("findings.Export: Cannot read evidence \"%s\".", evidence_id)
		}

		data = buffer
		default_path = filepath.Join("work", "evidence", id, evidence_id+"-"+filepath.Base(source))
		content_type = evidence.ContentType

	} else if export_type == "report" {

		report, _, err0 := buildFindingsReport(tool, "", false, "", false)

		if err0 != nil {
			return "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

		data = []byte(report)
		default_path = filepath.Join("work", "findings", "report.md")
		content_type = "text/markdown"

	} else if export_type == "scan" {

		if has_id == false || types.IsScanID(id) == false {
			return "", fmt.Errorf("findings.Export: Invalid parameter \"id\".")
		}

		buffer, err0 := utils_findings.ReadScanArtifact(tool.Playground, id)

		if err0 != nil {
			return "", fmt.Errorf("findings.Export: %s", err0.Error())
		}

		data = buffer
		default_path = filepath.Join("work", "scans", id+"-results.jsonl")
		content_type = "application/x-ndjson"

	} else {
		return "", fmt.Errorf("findings.Export: Invalid parameter \"type\".")
	}

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

	resolved, err1 := utils_sandbox.ResolvePath(sandbox, default_path)

	if err1 != nil {
		return "", fmt.Errorf("findings.Export: Cannot resolve export path.")
	}

	if err2 := os.WriteFile(resolved, data, 0644); err2 != nil {
		return "", fmt.Errorf("findings.Export: Cannot write export file.")
	}

	display, err3 := utils_sandbox.SanitizePath(sandbox, resolved)

	if err3 != nil {
		return "", fmt.Errorf("findings.Export: Cannot resolve export path.")
	}

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
