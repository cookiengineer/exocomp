package tools

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
		note,   ok3 := arguments["note"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Fix(path, method, note)
		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"note\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"method\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "FixAll" {

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
			return tool.Fix(path, method)
		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"method\" is not a string.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameter \"path\" is not a string.")
		} else {
			return "", fmt.Errorf("bugs.%s: %s", method, "Invalid parameters.")
		}

	}

}

func (tool *Bugs) Add(path string, method string, note string) (string, error) {
}

func (tool *Bugs) Fix(path string, method string, note string) (string, error) {
}

func (tool *Bugs) FixAll(path string, method string) (string, error) {
}

func (tool *Bugs) Search(path string, method string) (string, error) {
}
