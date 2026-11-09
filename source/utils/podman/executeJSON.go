package podman

import "bytes"
import "context"
import "os/exec"
import "strings"
import "time"

func executeJSON(arguments []string, timeout time.Duration, name string, image string) (string, error) {

	podman, err0 := podmanPath()

	if err0 == nil {

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		stdout := bytes.NewBuffer(make([]byte, 0, 4096))
		stderr := bytes.NewBuffer(make([]byte, 0, 4096))

		defer cancel()

		cmd := exec.CommandContext(ctx, podman, arguments...)
		cmd.Stdin = strings.NewReader("")
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err1 := cmd.Run()

		if ctx.Err() == context.DeadlineExceeded {

			return "", sanitizeError(name, image, stderr.String()+stdout.String(), errTimeout)

		} else if err1 == nil {

			return stdout.String(), nil

		} else {
			return "", sanitizeError(name, image, stderr.String()+stdout.String(), err1)
		}

	} else {
		return "", sanitizeError(name, image, "", err0)
	}

}
