package podman

import "fmt"
import "path/filepath"
import "time"

func Create(image string, name string, command string, arguments []string, network string, sandbox string, label string) error {

	absolute, err0 := filepath.Abs(sandbox)

	if err0 == nil {

		podman_arguments := []string{
			"run",
			"-d",
			"--pull=missing",
			"--name", name,
			"--label", managedKey + "=true",
			"--label", labelKey + "=" + label,
			"--network", network,
			"--workdir", "/workspace",
			"--volume", absolute + ":/workspace:Z",
			"--tmpfs", "/workspace/.exocomp:ro,size=1m",
			image,
			command,
		}

		podman_arguments = append(podman_arguments, arguments...)

		_, err1 := execute(podman_arguments, 10*time.Minute, 0, name, image)

		if err1 == nil {
			return nil
		} else {
			return err1
		}

	} else {
		return fmt.Errorf("Cannot resolve sandbox path.")
	}

}
