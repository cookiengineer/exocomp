package findings

import utils_findings "exocomp/utils/findings"

func readFindingsIndex(tool *Findings) error {

	summaries, err0 := utils_findings.ReadIndex(tool.Playground)

	if err0 == nil {

		for id, _ := range tool.contents {
			delete(tool.contents, id)
		}

		for _, summary := range summaries {
			tool.contents[summary.ID] = summary
		}

		return nil

	} else {
		return err0
	}

}
