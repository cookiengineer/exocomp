package podman

import "time"

func Start(name string) error {

	_, err := execute([]string{"start", name}, 30*time.Second, 0, name, "")

	return err

}
