package gadgets

import "fmt"

type Tasks struct {
	Sandbox string
	Tasks   map[string]bool
}

func NewTasks(sandbox string) *Tasks {

	tasks := &Tasks{
		Sandbox: sandbox,
		Tasks:   make(map[string]bool),
	}

	readTasks(tasks)

	return tasks

}

func (gadget *Tasks) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!gadget:tasks.List \"keyword\"",
		"",
		"#!gadget:tasks.Create \"The task description\"",
		"",
		"#!gadget:tasks.Finish \"The task description\"",
	}, "\n"), nil

}

func (gadget *Tasks) Create(arguments []string) (string, error) {

	if len(arguments) == 1 {

		label := strings.TrimSpace(arguments[0])
		found := ""

		for _, other := range tasks.Tasks {

			if strings.Contains(strings.ToLower(description), strings.ToLower(label)) == true {
				found = true
				break
			}

		}

		if found == "" {

			tasks.Tasks = append(tasks.Tasks, label)

			err := writeTasks(tasks)

			if err == nil {

				result := strings.Join([]string{
					fmt.Sprintf("#!tasks.Create: %s", label),
					fmt.Sprint("New Task created."),
				}, "\n")

				return result, nil

			} else {
				return "", err
			}

		} else {

			result := strings.Join([]string{
				fmt.Sprintf("#!tasks.Create: %s", label),
				fmt.Sprintf("Similar Task \"%s\" already exists.", found),
			}, "\n")

			return result, nil

		}

	} else {
		return "", fmt.Errorf("#!tasks.Create: Only one argument allowed")
	}

}

func (gadget *Tasks) Finish(arguments []string) (string, error) {
}

func (gadget *Tasks) List(arguments []string) (string, error) {

	if len(arguments) == 1 {

		keyword := strings.TrimSpace(arguments[0])
		found   := make([]string, 0)

		for _, description := range tasks.Tasks {

			if strings.Contains(strings.ToLower(description), strings.ToLower(keyword)) == true {
				found = append(found, description)
			}

		}

		if len(found) > 0 {

			result := make([]string, 0)
			result = append(result, fmt.Sprintf("#!tasks.List: %s", label))

			for f := 0; f < len(found); f++ {
			}

		} else {
			return "", fmt.Errorf("#!tasks.List: No matching tasks found")
		}

	} else {

		// TODO: List all

	}

}

