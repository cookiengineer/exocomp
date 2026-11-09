package podman

import "encoding/json"
import "time"

func Inspect(name string, label string) (*InspectSchema, error) {

	output, err1 := executeJSON([]string{"inspect", "--format", "json", name}, 30*time.Second, name, "")

	if err1 == nil {

		containers := make([]InspectSchema, 0)
		err2 := json.Unmarshal([]byte(output), &containers)

		if err2 == nil && len(containers) > 0 {

			if containers[0].Config.Labels[labelKey] == label {
				return &containers[0], nil
			} else {
				return nil, sanitizeError(name, "", output, errNotOwned)
			}

		} else {
			return nil, sanitizeError(name, "", output, errNotOwned)
		}

	} else {
		return nil, err1
	}

}
