package findings

import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_findings "exocomp/utils/findings"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
import "sync"

const (
	findingsDefaultLimit     = 50
	findingsMaxLimit         = 200
	findingsMaxTitle         = 200
	findingsMaxIdentifier    = 256
	findingsMaxEndpoint      = 2048
	findingsMaxText          = 16 * 1024
	findingsMaxEvidenceText  = 64 * 1024
	findingsMaxEvidenceCount = 32
	findingsMaxArtifactSize  = 4 * 1024 * 1024
)

type Findings struct {
	Methods    []string
	Playground string
	Sandbox    string
	AgentName  string
	mutex      *sync.RWMutex
	contents   map[string]types.FindingSummary
}

func NewFindings(methods []string, playground string, sandbox string, agent_name string) *Findings {

	tool := &Findings{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		AgentName:  agent_name,
		mutex:      &sync.RWMutex{},
		contents:   make(map[string]types.FindingSummary),
	}

	readFindingsIndex(tool)

	return tool

}

func (tool *Findings) Name() string {
	return "findings"
}

func (tool *Findings) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		parsed := utils_chat.Arguments(arguments)

		if method == "Add" {

			return tool.Add(parsed)

		} else if method == "Export" {

			return tool.Export(parsed)

		} else if method == "Import" {

			return tool.Import(parsed)

		} else if method == "BindEvidence" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.BindEvidence(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Get" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Get(id)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "List" {

			return tool.List(parsed)

		} else if method == "Report" {

			return tool.Report(parsed)

		} else if method == "Retest" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Retest(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Search" {

			return tool.Search(parsed)

		} else if method == "Update" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Update(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else if method == "Verify" {

			id, ok := parsed.String("id")

			if ok == true {
				return tool.Verify(id, parsed)
			} else {
				return "", fmt.Errorf("findings.%s: Invalid parameter \"id\" is not a string.", method)
			}

		} else {
			return "", fmt.Errorf("findings.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("findings.%s: Method not allowed.", method)
	}

}

func (tool *Findings) GetContent(id string) (any, error) {

	finding, err0 := utils_findings.ReadFinding(tool.Playground, id)

	if err0 == nil {
		return finding, nil
	} else {
		return nil, fmt.Errorf("findings.GetContent: %s", err0.Error())
	}

}

func (tool *Findings) GetContentIdentifiers() []string {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	readFindingsIndex(tool)

	result := make([]string, 0, len(tool.contents))

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	sort.Strings(result)

	return result

}

func (tool *Findings) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Findings) MarshalJSON() ([]byte, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	err0 := readFindingsIndex(tool)

	if err0 == nil {

		summaries := make([]types.FindingSummary, 0, len(tool.contents))

		for _, summary := range tool.contents {
			summaries = append(summaries, summary)
		}

		utils_findings.SortSummaries(summaries)

		return json.MarshalIndent(summaries, "", "\t")

	} else {
		return nil, fmt.Errorf("findings.MarshalJSON: %s", err0.Error())
	}

}

func (tool *Findings) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range FindingsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
