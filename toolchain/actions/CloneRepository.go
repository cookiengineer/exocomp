package actions

import "os"
import "os/exec"

func CloneRepository(cwd string, url string, folder string) error {

	stat, err0 := os.Stat(folder)

	if err0 == nil && stat.IsDir() {

		cmd := exec.Command(
			"git",
			"pull",
			"origin",
			"HEAD",
		)
		cmd.Dir = folder

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()

	} else if err0 != nil {

		cmd := exec.Command(
			"git",
			"clone",
			"--single-branch",
			"--depth", "1",
			url,
			folder,
		)
		cmd.Dir = cwd

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()

	}

	return nil

}
