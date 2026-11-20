package findings

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"

func importFindingsResults(tool *Findings, parsed []nucleiImportResult, invalid int, command string) (int, int, int, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err0 := readFindingsIndex(tool); err0 == nil {

		existing, err1 := utils_findings.ReadAll(tool.Playground)

		if err1 == nil {

			seen := make(map[string]bool, len(existing))

			for _, finding := range existing {
				seen[importDedupKey(finding)] = true
			}

			imported := 0
			skipped := 0

			for _, result := range parsed {

				if seen[importDedupKey(result.Finding)] == true {

					skipped++
					continue

				}

				if result.Request != "" {

					evidence := utils_chat.Arguments{
						"kind":    string(types.FindingEvidenceKindHTTPRequest),
						"title":   "Nuclei request",
						"source":  "nuclei",
						"content": result.Request,
					}

					if command != "" {
						evidence["command"] = command
					}

					if err2 := bindEvidenceToFinding(tool, &result.Finding, evidence); err2 != nil {
						invalid++
						continue
					}

				}

				if result.Response != "" {

					evidence := utils_chat.Arguments{
						"kind":    string(types.FindingEvidenceKindHTTPResponse),
						"title":   "Nuclei response",
						"source":  "nuclei",
						"content": result.Response,
					}

					if command != "" {
						evidence["command"] = command
					}

					if err3 := bindEvidenceToFinding(tool, &result.Finding, evidence); err3 != nil {
						invalid++
						continue
					}

				}

				if result.Finding.Validate() == true {

					if err4 := writeFinding(tool, result.Finding); err4 == nil {
						seen[importDedupKey(result.Finding)] = true
						imported++
					} else {
						return 0, 0, 0, err4
					}

				} else {
					invalid++
				}

			}

			return imported, skipped, invalid, nil

		} else {
			return 0, 0, 0, err1
		}

	} else {
		return 0, 0, 0, err0
	}

}
