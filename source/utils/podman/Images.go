package podman

import "encoding/json"
import "fmt"
import "time"

func Images() ([]ImagesSchema, error) {

	output, err1 := executeJSON([]string{"images", "--format", "json"}, 30*time.Second, "", "")

	if err1 == nil {

		images := make([]ImagesSchema, 0)
		err2 := json.Unmarshal([]byte(output), &images)

		if err2 == nil {
			return images, nil
		} else {
			return nil, fmt.Errorf("Cannot parse image list.")
		}

	} else {
		return nil, err1
	}

}
