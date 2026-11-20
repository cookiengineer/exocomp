package findings

import utils_chat "exocomp/utils/chat"
import utils_http "exocomp/utils/http"
import utils_sandbox "exocomp/utils/sandbox"
import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "strings"

func writeFindingsExport(tool *Findings, export_type string, data []byte, default_path string, content_type string, arguments utils_chat.Arguments) (string, error) {

	path, has_path := arguments.String("path")

	if has_path == true && strings.TrimSpace(path) != "" {
		default_path = path
	}

	sandbox := tool.Sandbox

	if sandbox == "" {
		sandbox = tool.Playground
	}

	if utils_sandbox.IsReservedPath(sandbox, default_path) == true {
		return "", fmt.Errorf("findings.Export: Path is reserved for internal state.")
	}

	resolved, display, err0 := writeFindingsExportFile(sandbox, default_path, data)

	if err0 == nil {

		if content_type == "" {
			content_type = utils_http.ContentTypeForPath(resolved)
		}

		sum := sha256.Sum256(data)

		return fmt.Sprintf(
			"findings.Export: type=%s path=\"%s\" sha256=%s size=%d content_type=%s",
			export_type,
			display,
			hex.EncodeToString(sum[:]),
			len(data),
			content_type,
		), nil

	} else {
		return "", fmt.Errorf("findings.Export: %s", err0.Error())
	}

}
