package build

import "os"
import "os/exec"

func runGit(cwd string, arguments ...string) error {

	cmd := exec.Command("git", arguments...)
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
