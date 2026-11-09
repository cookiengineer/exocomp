package podman

import "errors"
import "fmt"
import "io/fs"
import "os/exec"
import "strings"

func sanitizeError(name string, image string, output string, err error) error {

	result := strings.ToLower(output)

	if errors.Is(err, errTruncated) == true {
		return fmt.Errorf("Output exceeded 16MB limit.")
	} else if errors.Is(err, errTimeout) == true || errors.Is(err, errIdle) == true {
		return fmt.Errorf("Command timed out.")
	} else if errors.Is(err, exec.ErrNotFound) == true || errors.Is(err, fs.ErrNotExist) == true || strings.Contains(err.Error(), "executable file not found") == true {
		return fmt.Errorf("Program \"podman\" doesn't exist.")
	} else if errors.Is(err, errNotOwned) == true {
		return fmt.Errorf("Container \"%s\" does not exist.", name)
	} else if errors.Is(err, fs.ErrPermission) == true || strings.Contains(result, "permission denied") == true || strings.Contains(result, "operation not permitted") == true {

		if name == "" {
			return fmt.Errorf("Permission denied.")
		} else {
			return fmt.Errorf("Container \"%s\": Permission denied.", name)
		}

	} else if strings.Contains(result, "no such image") == true || strings.Contains(result, "unable to find image") == true || strings.Contains(result, "manifest unknown") == true {

		if image == "" {
			return fmt.Errorf("Image does not exist.")
		} else {
			return fmt.Errorf("Image \"%s\" does not exist.", image)
		}

	} else if strings.Contains(result, "no such container") == true || strings.Contains(result, "no such object") == true || strings.Contains(result, "no such name") == true {

		if name == "" {
			return fmt.Errorf("Container does not exist.")
		} else {
			return fmt.Errorf("Container \"%s\" does not exist.", name)
		}

	} else if strings.Contains(result, "already in use") == true || strings.Contains(result, "name is in use") == true {

		if name == "" {
			return fmt.Errorf("Container is already in use.")
		} else {
			return fmt.Errorf("Container \"%s\" is already in use.", name)
		}

	} else if name == "" {
		return fmt.Errorf("Command failed.")
	} else {
		return fmt.Errorf("Container \"%s\" command failed.", name)
	}

}
