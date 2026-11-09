package podman

import "time"

func Remove(name string) error {

	_, err := execute([]string{"rm", "-f", name}, 30*time.Second, 0, name, "")

	return err

}
