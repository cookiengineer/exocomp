package tools

import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "sort"
import "strings"

type bug_specification struct {
	IsFixed     bool   `json:"is_fixed"`
	File        string `json:"file"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

type Bugs struct {
	Playground string
	Sandbox    string
	contents   map[string]map[string]bug_specification // map[path][symbol]
}

func NewBugs(playground string, sandbox string) *Bugs {

	bugs := &Bugs{
		Playground: playground,
		Sandbox:    sandbox,
		contents:   make(map[string]map[string]bug_specification),
	}

	readBugs(bugs)

	return bugs

}

func (tool *Bugs) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "Add" {

		path,        ok1 := arguments["path"].(string)
		symbol,      ok2 := arguments["symbol"].(string)
		description, ok3 := arguments["description"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Add(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol), utils_fmt.FormatSingleLine(description))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"description\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Fix" {

		path,   ok1 := arguments["path"].(string)
		symbol, ok2 := arguments["symbol"].(string)

		if ok1 == true && ok2 == true {
			return tool.Fix(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"symbol\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Search" {

		path,   ok1 := arguments["path"].(string)
		symbol, ok2 := arguments["symbol"].(string)

		if ok1 == true && ok2 == true {
			return tool.Search(utils_fmt.FormatFilePath(path), utils_fmt.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return tool.Search(utils_fmt.FormatFilePath(path), "")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else {
		return "", fmt.Errorf("bugs.%s: Invalid method.", method)
	}

}

func (tool *Bugs) Add(path string, symbol string, description string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			_, ok1 := tool.contents[internal_path]

			if ok1 == false {
				tool.contents[internal_path] = make(map[string]bug_specification)
			}

			_, ok2 := tool.contents[internal_path][symbol]

			if ok2 == false {

				tool.contents[internal_path][symbol] = bug_specification{
					IsFixed:     false,
					File:        internal_path,
					Symbol:      symbol,
					Description: description,
				}

				err3 := writeBugs(tool)

				if err3 == nil {
					return fmt.Sprintf("bugs.Add: Bug report with %d B written.", len(description)), nil
				} else {
					return "", fmt.Errorf("bugs.Add: %s", err3.Error())
				}

			} else {

				bug_report := tool.contents[internal_path][symbol]
				bug_report.Description = description
				bug_report.IsFixed     = false
				tool.contents[internal_path][symbol] = bug_report

				err3 := writeBugs(tool)

				if err3 == nil {
					return fmt.Sprintf("bugs.Add: Bug report with %d B updated.", len(description)), nil
				} else {
					return "", fmt.Errorf("bugs.Add: %s", err3.Error())
				}

			}

		} else {
			return "", fmt.Errorf("bugs.Add: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("bugs.Add: %s", err1.Error())
	}

}

func (tool *Bugs) Fix(path string, symbol string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			_, ok1 := tool.contents[internal_path]

			if ok1 == true {

				bug_report, ok2 := tool.contents[internal_path][symbol]

				if ok2 == true {

					bug_report.IsFixed = true
					tool.contents[internal_path][symbol] = bug_report

					err3 := writeBugs(tool)

					if err3 == nil {
						return fmt.Sprintf("bugs.Fix: Bug report marked as fixed."), nil
					} else {
						return "", fmt.Errorf("bugs.Fix: %s", err3.Error())
					}

				} else {
					return "", fmt.Errorf("bugs.Fix: No bug report available for path \"%s\" and symbol \"%s\"", path, symbol)
				}

			} else {
				return "", fmt.Errorf("bugs.Fix: No bug reports available for path \"%s\".", path)
			}

		} else {
			return "", fmt.Errorf("bugs.Fix: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("bugs.Fix: %s", err1.Error())
	}

}

func (tool *Bugs) List() (string, error) {

	lines := make([]string, 0)

	for _, bug_reports := range tool.contents {

		for _, bug_report := range bug_reports {

			if bug_report.IsFixed == false {

				resolved_path, err1 := resolveSandboxPath(tool.Playground, bug_report.File)

				if err1 == nil {

					sandbox_path, err2 := sanitizeSandboxPath(tool.Sandbox, resolved_path)

					if err2 == nil {
						lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Description: %s", sandbox_path, bug_report.Symbol, bug_report.Description))
					}

				}

			}

		}

	}

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("bugs.List: %d unfixed bug reports.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Bugs) Search(path string, symbol string) (string, error) {

	tmp1, err1 := resolveSandboxPath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := sanitizeSandboxPath(tool.Playground, tmp1)

		if err2 == nil {

			if symbol != "" {

				lines  := make([]string, 0)
				_, ok1 := tool.contents[internal_path]

				if ok1 == true {

					bug_report, ok2 := tool.contents[internal_path][symbol]

					if ok2 == true {

						if bug_report.IsFixed == false {

							sandbox_path, err3 := sanitizeSandboxPath(tool.Sandbox, bug_report.File)

							if err3 == nil {
								lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Description: %s", sandbox_path, bug_report.Symbol, bug_report.Description))
							}

						}

					}

				}

				sort.Strings(lines)

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("bugs.Search: %s#%s contains %d bug reports.", path, symbol, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			} else {

				lines            := make([]string, 0)
				bug_reports, ok1 := tool.contents[internal_path]

				if ok1 == true {

					for _, bug_report := range bug_reports {

						if bug_report.IsFixed == false {

							sandbox_path, err3 := sanitizeSandboxPath(tool.Sandbox, bug_report.File)

							if err3 == nil {
								lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Description: %s", sandbox_path, bug_report.Symbol, bug_report.Description))
							}

						}

					}

				}

				sort.Strings(lines)

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("bugs.Search: %s contains %d bug reports.", path, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			}

		} else {
			return "", fmt.Errorf("bugs.Search: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("bugs.Search: %s", err1.Error())
	}

}
