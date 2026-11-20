package assets

import "exocomp/types"
import utils_http "exocomp/utils/http"
import utils_sandbox "exocomp/utils/sandbox"
import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "path/filepath"
import "strings"
import "time"

func runAssetExport(tool *Assets, matches []types.Asset, format string, path string, has_path bool) (string, error) {

	if has_path == false || strings.TrimSpace(path) == "" {

		extension := ".txt"

		if format == "json" {
			extension = ".json"
		}

		path = filepath.Join("work", fmt.Sprintf("assets-%s%s", time.Now().Format("20060102-150405"), extension))

	}

	sandbox := tool.Sandbox

	if sandbox == "" {
		sandbox = tool.Playground
	}

	if utils_sandbox.IsReservedPath(sandbox, path) == true {
		return "", fmt.Errorf("assets.Export: Path is reserved for internal state.")
	}

	data, err0 := formatAssets(matches, format)

	if err0 == nil {

		resolved, display, err1 := writeAssetExportFile(sandbox, path, data)

		if err1 == nil {

			sum := sha256.Sum256(data)
			content_type := utils_http.ContentTypeForPath(resolved)

			return fmt.Sprintf(
				"assets.Export: format=%s path=\"%s\" count=%d sha256=%s size=%d content_type=%s",
				format,
				display,
				len(matches),
				hex.EncodeToString(sum[:]),
				len(data),
				content_type,
			), nil

		} else {
			return "", fmt.Errorf("assets.Export: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Export: %s", err0.Error())
	}

}
