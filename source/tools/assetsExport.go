package tools

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_http "exocomp/utils/http"
import utils_sandbox "exocomp/utils/sandbox"
import "crypto/sha256"
import "encoding/hex"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"
import "strings"
import "time"

func (tool *Assets) Export(arguments utils_chat.Arguments) (string, error) {

	format := "hosts"

	if value, ok := arguments.String("format"); ok == true {
		format = strings.ToLower(strings.TrimSpace(value))
	}

	if format != "hosts" && format != "urls" && format != "json" {
		return "", fmt.Errorf("assets.Export: Invalid parameter \"format\".")
	}

	type_filter, has_type := arguments.String("type")
	query_filter, has_query := arguments.String("query")
	in_scope_filter, has_in_scope := arguments.Bool("in_scope")
	tested_filter, has_tested := arguments.Bool("tested")
	path, has_path := arguments.String("path")
	limit := assetExportLimit(arguments)

	if has_type == true && types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == false {
		return "", fmt.Errorf("assets.Export: Invalid parameter \"type\".")
	}

	matches, err0 := collectExportAssets(tool, type_filter, has_type, query_filter, has_query, in_scope_filter, has_in_scope, tested_filter, has_tested, limit)

	if err0 == nil {
		return runAssetExport(tool, matches, format, path, has_path)
	} else {
		return "", fmt.Errorf("assets.Export: %s", err0.Error())
	}

}

func collectExportAssets(tool *Assets, type_filter string, has_type bool, query_filter string, has_query bool, in_scope_filter bool, has_in_scope bool, tested_filter bool, has_tested bool, limit int) ([]types.Asset, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		matches := make([]types.Asset, 0)

		for _, asset := range tool.contents {

			if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
				continue
			} else if has_query == true && assetMatchesQuery(asset, query_filter) == false {
				continue
			} else if has_in_scope == true && asset.InScope != in_scope_filter {
				continue
			} else if has_tested == true && asset.Tested != tested_filter {
				continue
			}

			matches = append(matches, asset)

		}

		sortAssets(matches)

		if len(matches) > limit {
			matches = matches[0:limit]
		}

		return matches, nil

	} else {
		return nil, err
	}

}

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

func writeAssetExportFile(sandbox string, path string, data []byte) (string, string, error) {

	resolved, err0 := utils_sandbox.ResolvePath(sandbox, path)

	if err0 == nil {

		if err1 := os.WriteFile(resolved, data, 0644); err1 == nil {

			display, err2 := utils_sandbox.SanitizePath(sandbox, resolved)

			if err2 == nil {
				return resolved, display, nil
			} else {
				return "", "", fmt.Errorf("Cannot resolve export path.")
			}

		} else {
			return "", "", fmt.Errorf("Cannot write export file.")
		}

	} else {
		return "", "", fmt.Errorf("Cannot resolve export path.")
	}

}

func assetExportLimit(arguments utils_chat.Arguments) int {

	limit := 1000

	if value, ok := arguments.Int("limit"); ok == true && value > 0 {
		limit = value
	}

	if limit > 10000 {
		limit = 10000
	}

	return limit

}

func formatAssets(entries []types.Asset, format string) ([]byte, error) {

	if format == "json" {
		return json.MarshalIndent(entries, "", "\t")
	}

	lines := make([]string, 0, len(entries))

	for _, entry := range entries {

		value := entry.Value

		if format == "urls" {

			if entry.URL != "" {
				value = entry.URL
			} else if strings.HasPrefix(entry.Value, "http://") == false && strings.HasPrefix(entry.Value, "https://") == false {
				continue
			}

		}

		if strings.TrimSpace(value) != "" {
			lines = append(lines, value)
		}

	}

	return []byte(strings.Join(lines, "\n") + "\n"), nil

}
