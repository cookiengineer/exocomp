package assets

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
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

func (tool *Assets) Call(method string, arguments map[string]interface{}) (string, error) {

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
