package tools

import "exocomp/types"
import utils_findings "exocomp/utils/findings"

func writeFinding(tool *Findings, finding types.Finding) error {

	err0 := utils_findings.WriteFinding(tool.Playground, finding)

	if err0 == nil {

		tool.contents[finding.ID] = types.FindingSummaryFromFinding(finding)

		return writeFindingsIndex(tool)

	} else {
		return err0
	}

}
