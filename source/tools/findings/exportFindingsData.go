package findings

import utils_chat "exocomp/utils/chat"
import "fmt"

func exportFindingsData(tool *Findings, export_type string, arguments utils_chat.Arguments) ([]byte, string, string, error) {

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
