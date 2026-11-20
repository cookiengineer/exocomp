package assets

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
import "strings"
import "sync"

const (
	assetsDefaultLimit = 50
	assetsMaxLimit     = 200
	assetsMaxValue     = 2048
	assetsMaxText      = 1024
)

type Assets struct {
	Methods    []string
	Playground string
	Sandbox    string
	mutex      *sync.RWMutex
	contents   map[string]types.Asset
}

func NewAssets(methods []string, playground string, sandbox string) *Assets {

	tool := &Assets{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		mutex:      &sync.RWMutex{},
		contents:   make(map[string]types.Asset),
	}

	readAssetsIndex(tool)

	return tool

}

func (tool *Assets) Name() string {
	return "assets"
}

func (tool *Assets) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		parsed := utils_chat.Arguments(arguments)

		if method == "Add" {

			return tool.Add(parsed)

		} else if method == "Export" {

			return tool.Export(parsed)

		} else if method == "Import" {

			return tool.Import(parsed)

		} else if method == "Get" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Get(id)
			} else {
				return "", fmt.Errorf("assets.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "List" {

			return tool.List(parsed)

		} else if method == "Search" {

			return tool.Search(parsed)

		} else if method == "Update" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Update(id, parsed)
			} else {
				return "", fmt.Errorf("assets.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Neighbors" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Neighbors(id, parsed)
			} else {
				return "", fmt.Errorf("assets.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Scope" {

			return tool.Scope(parsed)

		} else if method == "Coverage" {

			return tool.Coverage()

		} else {
			return "", fmt.Errorf("assets.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("assets.%s: Method not allowed.", method)
	}

}

func (tool *Assets) GetContent(id string) (any, error) {

	if types.IsAssetID(id) == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		asset, err0 := getAsset(tool, id)

		if err0 == nil {
			return asset, nil
		} else {
			return nil, fmt.Errorf("assets.GetContent: %s", err0.Error())
		}

	} else {
		return nil, fmt.Errorf("assets.GetContent: Invalid asset ID \"%s\".", id)
	}

}

func (tool *Assets) GetContentIdentifiers() []string {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	readAssetsIndex(tool)

	result := make([]string, 0, len(tool.contents))

	for id := range tool.contents {
		result = append(result, id)
	}

	sort.Strings(result)

	return result

}

func (tool *Assets) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Assets) MarshalJSON() ([]byte, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		entries := make([]types.Asset, 0, len(tool.contents))

		for _, asset := range tool.contents {
			entries = append(entries, asset)
		}

		sortAssets(entries)

		return json.MarshalIndent(entries, "", "\t")

	} else {
		return nil, fmt.Errorf("assets.MarshalJSON: %s", err.Error())
	}

}

func (tool *Assets) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range AssetsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}

func (tool *Assets) Add(arguments utils_chat.Arguments) (string, error) {

	type_value, has_type := arguments.String("type")
	value, has_value := arguments.String("value")

	if has_type == true && has_value == true {

		asset := types.Asset{
			Type:       types.AssetType(strings.ToLower(strings.TrimSpace(type_value))),
			Value:      value,
			Technology: make([]string, 0),
			Props:      make(map[string]string),
			Tags:       make([]string, 0),
		}

		if asset.Type.IsValid() == true {

			if parent_id, ok := arguments.String("parent_id"); ok == true {
				asset.ParentID = parent_id
			}

			if domain, ok := arguments.String("domain"); ok == true {
				asset.Domain = utils_findings.LimitText(utils_findings.RedactSecrets(domain), assetsMaxValue)
			}

			if ip, ok := arguments.String("ip"); ok == true {
				asset.IP = utils_findings.LimitText(utils_findings.RedactSecrets(ip), assetsMaxValue)
			}

			if port, ok := arguments.Int("port"); ok == true {
				asset.Port = port
			}

			if raw_url, ok := arguments.String("url"); ok == true {
				asset.URL = utils_findings.LimitText(utils_findings.RedactSecrets(raw_url), assetsMaxValue)
			}

			if endpoint, ok := arguments.String("endpoint"); ok == true {
				asset.Endpoint = utils_findings.LimitText(utils_findings.RedactSecrets(endpoint), assetsMaxValue)
			}

			if technology, ok := arguments.Strings("technology"); ok == true {
				asset.Technology = normalizeStringList(technology, assetsMaxText)
			}

			if source, ok := arguments.String("source"); ok == true {
				asset.Source = utils_findings.LimitText(utils_findings.RedactSecrets(source), assetsMaxText)
			}

			if confidence, ok := arguments.String("confidence"); ok == true {
				asset.Confidence = types.FindingConfidence(strings.ToLower(strings.TrimSpace(confidence)))
			}

			if in_scope, ok := arguments.Bool("in_scope"); ok == true {
				asset.InScope = in_scope
			}

			if tested, ok := arguments.Bool("tested"); ok == true {
				asset.Tested = tested
			}

			if tags, ok := arguments.Strings("tags"); ok == true {
				asset.Tags = normalizeStringList(tags, assetsMaxText)
			}

			if props, ok := arguments.Map("props"); ok == true {
				asset.Props = normalizeProps(props)
			}

			asset.Value = utils_findings.NormalizeAssetValue(asset.Type, asset.Value)
			asset.Value = utils_findings.LimitText(utils_findings.RedactSecrets(asset.Value), assetsMaxValue)

			if err := validateAsset(&asset); err == nil {

				stored, created, err := storeAsset(tool, asset)

				if err == nil {

					if created == true {
						return fmt.Sprintf("assets.Add: Asset \"%s\" created.", stored.ID), nil
					} else {
						return fmt.Sprintf("assets.Add: Asset \"%s\" updated.", stored.ID), nil
					}

				} else {
					return "", fmt.Errorf("assets.Add: %s", err.Error())
				}

			} else {
				return "", fmt.Errorf("assets.Add: %s", err.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Add: Invalid asset type \"%s\".", type_value)
		}

	} else {
		return "", fmt.Errorf("assets.Add: Invalid parameters \"type\" and \"value\" are required.")
	}

}

func (tool *Assets) Coverage() (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	if err := readAssetsIndex(tool); err == nil {

		counts := make(map[types.AssetType][3]int)

		for _, asset := range tool.contents {

			count := counts[asset.Type]

			if asset.InScope == true {
				count[0]++
			}

			if asset.Tested == true {
				count[1]++
			}

			count[2]++
			counts[asset.Type] = count

		}

		types_list := make([]string, 0, len(counts))

		for asset_type := range counts {
			types_list = append(types_list, string(asset_type))
		}

		sort.Strings(types_list)

		lines := make([]string, 0, len(types_list))

		for _, name := range types_list {

			count := counts[types.AssetType(name)]
			lines = append(lines, fmt.Sprintf("%s: %d assets, %d in scope, %d tested", name, count[2], count[0], count[1]))

		}

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("assets.Coverage: %d assets across %d types.", len(tool.contents), len(types_list)))
		result = append(result, lines...)

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("assets.Coverage: %s", err.Error())
	}

}

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

func (tool *Assets) Get(id string) (string, error) {

	if types.IsAssetID(id) == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		asset, err0 := getAsset(tool, id)

		if err0 == nil {

			buffer, err1 := json.MarshalIndent(asset, "", "\t")

			if err1 == nil {
				return strings.Join([]string{fmt.Sprintf("assets.Get: %s", id), string(buffer)}, "\n"), nil
			} else {
				return "", fmt.Errorf("assets.Get: Cannot serialize asset \"%s\".", id)
			}

		} else {
			return "", fmt.Errorf("assets.Get: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Get: Invalid asset ID \"%s\".", id)
	}

}

func (tool *Assets) Import(arguments utils_chat.Arguments) (string, error) {

	path, ok := arguments.String("path")

	if ok == true && strings.TrimSpace(path) != "" {

		format, ok := arguments.String("format")

		if ok == true {

			format = strings.ToLower(strings.TrimSpace(format))

			if isSupportedAssetFormat(format) == true {

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

			} else {
				return "", fmt.Errorf("assets.Import: Unsupported format \"%s\".", format)
			}

		} else {
			return "", fmt.Errorf("assets.Import: Invalid parameter \"format\" is not a string.")
		}

	} else {
		return "", fmt.Errorf("assets.Import: Invalid parameter \"path\" is not a non-empty string.")
	}

}

func (tool *Assets) List(arguments utils_chat.Arguments) (string, error) {

	type_filter, has_type := arguments.String("type")
	query, has_query := arguments.String("query")
	in_scope_filter, has_in_scope := arguments.Bool("in_scope")
	tested_filter, has_tested := arguments.Bool("tested")
	limit := assetLimit(arguments)

	if has_type == false || types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		if err := readAssetsIndex(tool); err == nil {

			matches := make([]types.Asset, 0)

			for _, asset := range tool.contents {

				if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
					continue
				} else if has_query == true && assetMatchesQuery(asset, query) == false {
					continue
				} else if has_in_scope == true && asset.InScope != in_scope_filter {
					continue
				} else if has_tested == true && asset.Tested != tested_filter {
					continue
				}

				matches = append(matches, asset)

			}

			sortAssets(matches)

			total := len(matches)

			if len(matches) > limit {
				matches = matches[0:limit]
			}

			lines := make([]string, 0, len(matches))

			for _, asset := range matches {
				lines = append(lines, formatAssetSummary(asset))
			}

			result := make([]string, 0)
			result = append(result, fmt.Sprintf("assets.List: %d assets (showing %d).", total, len(lines)))
			result = append(result, lines...)

			return strings.Join(result, "\n"), nil

		} else {
			return "", fmt.Errorf("assets.List: %s", err.Error())
		}

	} else {
		return "", fmt.Errorf("assets.List: Invalid parameter \"type\".")
	}

}

func (tool *Assets) Neighbors(id string, arguments utils_chat.Arguments) (string, error) {

	if types.IsAssetID(id) == true {

		direction := "both"

		if value, ok := arguments.String("direction"); ok == true {
			direction = strings.ToLower(strings.TrimSpace(value))
		}

		if direction == "parent" || direction == "children" || direction == "both" {

			limit := assetLimit(arguments)

			tool.mutex.Lock()
			defer tool.mutex.Unlock()

			asset, err0 := getAsset(tool, id)

			if err0 == nil {

				lines := make([]string, 0)

				if (direction == "parent" || direction == "both") && asset.ParentID != "" {

					if parent, ok := tool.contents[asset.ParentID]; ok == true {
						lines = append(lines, "Parent: "+formatAssetSummary(parent))
					}

				}

				if direction == "children" || direction == "both" {

					children := make([]types.Asset, 0)

					for _, candidate := range tool.contents {

						if candidate.ParentID == id {
							children = append(children, candidate)
						}

					}

					sortAssets(children)

					if len(children) > limit {
						children = children[0:limit]
					}

					if len(children) > 0 {
						lines = append(lines, fmt.Sprintf("Children: %d", len(children)))
					}

					for _, child := range children {
						lines = append(lines, formatAssetSummary(child))
					}

				}

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("assets.Neighbors: %s", id))

				if len(lines) == 0 {
					result = append(result, "No neighbors.")
				} else {
					result = append(result, lines...)
				}

				return strings.Join(result, "\n"), nil

			} else {
				return "", fmt.Errorf("assets.Neighbors: %s", err0.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Neighbors: Invalid parameter \"direction\".")
		}

	} else {
		return "", fmt.Errorf("assets.Neighbors: Invalid asset ID \"%s\".", id)
	}
}

func (tool *Assets) Scope(arguments utils_chat.Arguments) (string, error) {

	type_filter, has_type := arguments.String("type")
	value_filter, has_value := arguments.String("value")
	query, has_query := arguments.String("query")
	has_filter := has_type == true || has_value == true || has_query == true
	in_scope := true

	if provided, ok := arguments.Bool("in_scope"); ok == true {
		in_scope = provided
	}

	if has_type == false || types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		if err := readAssetsIndex(tool); err == nil {

			if has_filter == false {

				total := len(tool.contents)
				scoped := 0
				tested := 0

				for _, asset := range tool.contents {

					if asset.InScope == true {
						scoped++
					}

					if asset.Tested == true {
						tested++
					}

				}

				return fmt.Sprintf("assets.Scope: %d assets, %d in scope, %d tested.", total, scoped, tested), nil

			} else {

				updated := 0

				for id, asset := range tool.contents {

					if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
						continue
					} else if has_value == true && strings.EqualFold(asset.Value, strings.TrimSpace(value_filter)) == false {
						continue
					} else if has_query == true && assetMatchesQuery(asset, query) == false {
						continue
					}

					asset.InScope = in_scope
					asset.LastSeen = schemas.NewDatetime()
					tool.contents[id] = asset
					updated++

				}

				if err := writeAssetsIndex(tool); err == nil {
					return fmt.Sprintf("assets.Scope: %d assets marked in_scope=%t.", updated, in_scope), nil
				} else {
					return "", fmt.Errorf("assets.Scope: %s", err.Error())
				}

			}

		} else {
			return "", fmt.Errorf("assets.Scope: %s", err.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Scope: Invalid parameter \"type\".")
	}

}

func (tool *Assets) Search(arguments utils_chat.Arguments) (string, error) {

	query, ok := arguments.String("query")

	if ok == true && strings.TrimSpace(query) != "" {

		type_filter, has_type := arguments.String("type")
		limit := assetLimit(arguments)

		if has_type == false || types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))).IsValid() == true {

			tool.mutex.Lock()
			defer tool.mutex.Unlock()

			if err := readAssetsIndex(tool); err == nil {

				matches := make([]types.Asset, 0)

				for _, asset := range tool.contents {

					if has_type == true && asset.Type != types.AssetType(strings.ToLower(strings.TrimSpace(type_filter))) {
						continue
					} else if assetMatchesQuery(asset, query) == false {
						continue
					}

					matches = append(matches, asset)

				}

				sortAssets(matches)

				total := len(matches)

				if len(matches) > limit {
					matches = matches[0:limit]
				}

				lines := make([]string, 0, len(matches))

				for _, asset := range matches {
					lines = append(lines, formatAssetSummary(asset))
				}

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("assets.Search: %d matches (showing %d).", total, len(lines)))
				result = append(result, lines...)

				return strings.Join(result, "\n"), nil

			} else {
				return "", fmt.Errorf("assets.Search: %s", err.Error())
			}

		} else {
			return "", fmt.Errorf("assets.Search: Invalid parameter \"type\".")
		}

	} else {
		return "", fmt.Errorf("assets.Search: Invalid parameter \"query\" is not a non-empty string.")
	}

}

func (tool *Assets) Update(id string, arguments utils_chat.Arguments) (string, error) {

	if types.IsAssetID(id) == true {

		tool.mutex.Lock()
		defer tool.mutex.Unlock()

		asset, err0 := getAsset(tool, id)

		if err0 == nil {
			return updateAsset(tool, id, asset, arguments)
		} else {
			return "", fmt.Errorf("assets.Update: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("assets.Update: Invalid asset ID \"%s\".", id)
	}

}
