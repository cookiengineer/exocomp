package findings

import utils_findings "exocomp/utils/findings"
import "encoding/json"
import "fmt"
import "strings"

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
