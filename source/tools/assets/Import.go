package assets

import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "fmt"
import "strings"

func (tool *Assets) Import(arguments utils_chat.Arguments) (string, error) {

	path, ok := arguments.String("path")

	if ok == false || strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("assets.Import: Invalid parameter \"path\" is not a non-empty string.")
	}

	format, ok := arguments.String("format")

	if ok == false {
		return "", fmt.Errorf("assets.Import: Invalid parameter \"format\" is not a string.")
	}

	format = strings.ToLower(strings.TrimSpace(format))

	if isSupportedAssetFormat(format) == false {
		return "", fmt.Errorf("assets.Import: Unsupported format \"%s\".", format)
	}

	target, _ := arguments.String("target")
	in_scope, _ := arguments.Bool("in_scope")
	sandbox := tool.Sandbox

	if sandbox == "" {
		sandbox = tool.Playground
	}

	data, err0 := utils_findings.ReadScanArtifactFile(sandbox, path, "assets.Import")

	if err0 == nil {

		parsed, invalid, err1 := parseAssetResults(format, data, target, in_scope)

		if err1 == nil {

			created, updated, invalid, err2 := importAssetResults(tool, parsed, invalid)

			if err2 == nil {
				return fmt.Sprintf("assets.Import: Created %d assets, updated %d, skipped %d invalid results (format %s).", created, updated, invalid, format), nil
			} else {
				return "", fmt.Errorf("assets.Import: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Import: %s", err1.Error())
		}

	} else {
		return "", err0
	}

}
