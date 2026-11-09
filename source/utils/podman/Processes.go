package podman

import "encoding/json"
import "fmt"
import "time"

func Processes(label string) ([]ProcessesSchema, error) {

	arguments := []string{
		"ps",
		"-a",
		"--filter", "label=" + labelKey + "=" + label,
		"--format", "json",
	}

	output, err1 := executeJSON(arguments, 30*time.Second, "", "")

	if err1 == nil {

		processes := make([]ProcessesSchema, 0)
		err2 := json.Unmarshal([]byte(output), &processes)

		if err2 == nil {
			return processes, nil
		} else {
			return nil, fmt.Errorf("Cannot parse container list.")
		}

	} else {
		return nil, err1
	}

}
