package tools

import "exocomp/utils"
import "fmt"
import "sort"
import "strings"

type Bugs struct {
	Sandbox  string
	//       map[./path/to/File.go#Symbol]map[note]is_fixed
	contents map[string]map[string]bool
}

func NewBugs(agent string, sandbox string) *Bugs {

	bugs := &Bugs{
		Sandbox:  sandbox,
		contents: make(map[string]map[string]bool),
	}

	readBugs(bugs)

	return bugs

}

func (tool *Bugs) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "Add" {

		path,   ok1 := arguments["path"].(string)
		symbol, ok2 := arguments["symbol"].(string)
		note,   ok3 := arguments["note"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Add(utils.FormatFilePath(path), utils.FormatSymbol(symbol), utils.FormatSingleLine(note))
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"note\" is not a string.")
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
			return tool.Fix(utils.FormatFilePath(path), utils.FormatSymbol(symbol))
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
			return tool.Search(utils.FormatFilePath(path), utils.FormatSymbol(symbol))
		} else if ok1 == true && ok2 == false {
			return tool.Search(utils.FormatFilePath(path), "")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else {
		return "", fmt.Errorf("bugs.%s: Invalid method.", method)
	}

}

func (tool *Bugs) Add(path string, symbol string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		anchor := fmt.Sprintf("%s#%s", path, symbol)
		_, ok1 := tool.contents[anchor]

		if ok1 == false {
			tool.contents[anchor] = make(map[string]bool)
		}

		tool.contents[anchor][note] = false
		writeBugs(tool)

		result := strings.Join([]string{
			fmt.Sprintf("bugs.Add: Bug report with %d B written.", len(note)),
		}, "\n")

		return result, nil

	} else {
		return "", fmt.Errorf("bugs.Add: %s", err0.Error())
	}

}

func (tool *Bugs) Fix(path string, symbol string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		anchor := fmt.Sprintf("%s#%s", path, symbol)
		_, ok1 := tool.contents[anchor]

		if ok1 == true {

			count := 0

			for note, _ := range tool.contents[anchor] {
				tool.contents[anchor][note] = true
				count++
			}

			result := strings.Join([]string{
				fmt.Sprintf("bugs.Fix: %d bug reports marked as fixed.", count),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("bugs.Fix: 0 bug reports marked as fixed.")
		}

	} else {
		return "", fmt.Errorf("bugs.Fix: %s", err0.Error())
	}

}

func (tool *Bugs) List() (string, error) {

	lines := make([]string, 0)

	for anchor, notes := range tool.contents {

		for note, is_fixed := range notes {

			if is_fixed == false {
				lines = append(lines, fmt.Sprintf("- [ ] `%s`: %s", anchor, note))
			}

		}

	}

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("bugs.List: %d bug reports.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Bugs) Search(path string, symbol string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		lines := make([]string, 0)

		if symbol != "" {

			anchor     := fmt.Sprintf("%s#%s", path, symbol)
			notes, ok1 := tool.contents[anchor]

			if ok1 == true {

				for note, is_fixed := range notes {

					if is_fixed == false {
						lines = append(lines, fmt.Sprintf("- [ ] `%s`: %s", anchor, note))
					}

				}

			}

		} else {

			for anchor, notes := range tool.contents {

				if strings.HasPrefix(anchor, path + "#") {

					for note, is_fixed := range notes {

						if is_fixed == false {
							lines = append(lines, fmt.Sprintf("- [ ] `%s`: %s", anchor, note))
						}

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

	} else {
		return "", fmt.Errorf("bugs.Search: %s", err0.Error())
	}

}
