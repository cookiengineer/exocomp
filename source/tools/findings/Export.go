package findings

import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Findings) Export(arguments utils_chat.Arguments) (string, error) {

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
