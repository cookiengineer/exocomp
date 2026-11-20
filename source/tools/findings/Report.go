package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Findings) Report(arguments utils_chat.Arguments) (string, error) {

	severity_filter, has_severity := arguments.String("severity")
	status_filter, has_status := arguments.String("status")

	if has_severity == true && types.FindingSeverity(strings.ToLower(strings.TrimSpace(severity_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.Report: Invalid parameter \"severity\".")
	}

	if has_status == true && types.FindingStatus(strings.ToLower(strings.TrimSpace(status_filter))).IsValid() == false {
		return "", fmt.Errorf("findings.Report: Invalid parameter \"status\".")
	}

	report, count, err0 := buildFindingsReport(tool, severity_filter, has_severity, status_filter, has_status)

	if err0 == nil {
		return fmt.Sprintf("findings.Report: %d findings.\n%s", count, report), nil
	} else {
		return "", fmt.Errorf("findings.Report: %s", err0.Error())
	}

}
