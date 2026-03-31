package tools

import "exocomp/utils"
import "fmt"
import "sort"
import "strings"

type Bugs struct {
	Sandbox  string
	contents map[string]map[string]bool // map[File.go#Method]map[bug_description]is_fixed
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

	if method == "Add" {

		path,   ok1 := arguments["path"].(string)
		method, ok2 := arguments["method"].(string)
		note,   ok3 := arguments["note"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Add(path, method, note)
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"method\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Fix" {

		path,   ok1 := arguments["path"].(string)
		method, ok2 := arguments["method"].(string)

		if ok1 == true && ok2 == true {
			return tool.Fix(path, method)
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"method\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Search" {

		path,   ok1 := arguments["path"].(string)
		method, ok2 := arguments["method"].(string)

		if ok1 == true && ok2 == true {
			return tool.Search(path, method)
		} else if ok1 == true && ok2 == false {
			return tool.Search(path, "")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else {
		return "", fmt.Errorf("bugs.%s: Invalid method.", method)
	}

}

func (tool *Bugs) Add(path string, method string, note string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		anchor := fmt.Sprintf("%s#%s", path, method)
		_, ok1 := tool.contents[anchor]

		if ok1 == false {
			tool.contents[anchor] = make(map[string]bool)
		}

		message := utils.FormatSingleLine(note)
		tool.contents[anchor][message] = false
		writeBugs(tool)

		result := strings.Join([]string{
			fmt.Sprintf("bugs.Add: Report with %d B written.", len(message)),
		}, "\n")

		return result, nil

	} else {
		return "", fmt.Errorf("bugs.Add: %s", err0.Error())
	}

}

func (tool *Bugs) Fix(path string, method string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		anchor := fmt.Sprintf("%s#%s", path, method)
		_, ok1 := tool.contents[anchor]

		if ok1 == true {

			count := 0

			for message, _ := range tool.contents[anchor] {
				tool.contents[anchor][message] = true
				count++
			}

			result := strings.Join([]string{
				fmt.Sprintf("bugs.Fix: %d Reports marked as fixed.", count),
			}, "\n")

			return result, nil

		} else {
			return "", fmt.Errorf("bugs.Fix: No matching Bug Reports found.")
		}

	} else {
		return "", fmt.Errorf("bugs.Fix: %s", err0.Error())
	}

}

func (tool *Bugs) Search(path string, method string) (string, error) {

	_, err0 := resolveSandboxPath(tool.Sandbox, path)

	if err0 == nil {

		lines := make([]string, 0)

		if method != "" {

			anchor     := fmt.Sprintf("%s#%s", path, method)
			notes, ok1 := tool.contents[anchor]

			if ok1 == true {

				for note, is_fixed := range notes {

					if is_fixed == false {
						lines = append(lines, fmt.Sprintf("- [ ] `%s`: %s", anchor, note))
					}

				}

			} else {
				return "", fmt.Errorf("bugs.Search: No matching Bug Reports found.")
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
		result = append(result, fmt.Sprintf("bugs.Search: %s contains %d items.", path, len(lines)))

		for l := 0; l < len(lines); l++ {
			result = append(result, lines[l])
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("bugs.Search: %s", err0.Error())
	}

}
