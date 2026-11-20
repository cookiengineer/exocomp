package assets

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"
import "strings"

func (tool *Assets) Export(arguments utils_chat.Arguments) (string, error) {

	format := "hosts"

	if value, ok := arguments.String("format"); ok == true {
		format = strings.ToLower(strings.TrimSpace(value))
	}

	if format == "hosts" || format == "urls" || format == "json" {

		type_filter, has_type := arguments.String("type")
		query_filter, has_query := arguments.String("query")
		in_scope_filter, has_in_scope := arguments.Bool("in_scope")
		tested_filter, has_tested := arguments.Bool("tested")
		path, has_path := arguments.String("path")
		limit := assetExportLimit(arguments)

		if has_type == false || types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == true {

			matches, err0 := collectExportAssets(tool, type_filter, has_type, query_filter, has_query, in_scope_filter, has_in_scope, tested_filter, has_tested, limit)

			if err0 == nil {
				return runAssetExport(tool, matches, format, path, has_path)
			} else {
				return "", fmt.Errorf("assets.Export: %s", err0.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Export: Invalid parameter \"type\".")
		}
	} else {
		return "", fmt.Errorf("assets.Export: Invalid parameter \"format\".")
	}
}
