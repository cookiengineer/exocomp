package utils

import "os"
import "os/exec"

func BuildWASM(cwd string, source string, output string) error {

	args := []string{
		"build",
		"-o", output,
		source,
	}

	env := os.Environ()

	env = append(env, "CGO_ENABLED=0")
	env = append(env, "GOOS=js")
	env = append(env, "GOARCH=wasm")

	cmd := exec.Command("go", args...)
	cmd.Dir = cwd
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
