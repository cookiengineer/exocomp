package findings

import "exocomp/types"
import utils_findings "exocomp/utils/findings"

func writeFindingsIndex(tool *Findings) error {

	summaries := make([]types.FindingSummary, 0, len(tool.contents))

	for _, summary := range tool.contents {
		summaries = append(summaries, summary)
	}

	return utils_findings.WriteIndex(tool.Playground, summaries)

}
