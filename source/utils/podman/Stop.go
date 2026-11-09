package podman

import "time"

func Stop(name string) error {

	_, err := execute([]string{"stop", name}, 30*time.Second, 0, name, "")

	return err

}
