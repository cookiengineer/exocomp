package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Findings) Import(arguments utils_chat.Arguments) (string, error) {

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

}
