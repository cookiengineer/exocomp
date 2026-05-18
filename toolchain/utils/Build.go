package utils

import "fmt"
import "os"
import "os/exec"
import "strings"

func Build(cwd string, source string, output string, tags []string, operating_system string) error {

	if operating_system == "windows" && strings.HasSuffix(output, ".exe") == false {
		output = fmt.Sprintf("%s.exe", output)
	}

	args := []string{"build"}
	env  := os.Environ()

	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, " "))
	}

	args = append(args, "-o", output)
	args = append(args, source)

	env = append(env, "CGO_ENABLED=0")
	env = append(env, fmt.Sprintf("GOOS=%s", operating_system))
	env = append(env, fmt.Sprintf("GOARCH=%s", "amd64"))

	cmd := exec.Command("go", args...)
	cmd.Dir = cwd
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
