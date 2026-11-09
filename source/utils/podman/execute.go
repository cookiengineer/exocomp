package podman

import utils_bytes "exocomp/utils/bytes"
import "context"
import "os/exec"
import "strings"
import "time"

const maxOutput = 16 * 1024 * 1024

func execute(arguments []string, timeout time.Duration, idle_timeout time.Duration, name string, image string) (string, error) {

	podman, err0 := podmanPath()

	if err0 == nil {

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		buffer := utils_bytes.NewContextBuffer(maxOutput, cancel)

		defer cancel()

		if idle_timeout > 0 {

			go func() {

				ticker := time.NewTicker(10 * time.Second)

			BackgroundLoop:
				for {

					select {

					case <-ctx.Done():

						break BackgroundLoop

					case <-ticker.C:

						last_write := buffer.LastWrite()

						if time.Since(last_write) > idle_timeout {

							cancel()
							break BackgroundLoop

						}

					}

				}

				ticker.Stop()

			}()

		}

		cmd := exec.CommandContext(ctx, podman, arguments...)
		cmd.Stdin = strings.NewReader("")
		cmd.Stdout = buffer
		cmd.Stderr = buffer

		err1 := cmd.Run()

		if buffer.IsTruncated() == true {

			return "", sanitizeError(name, image, buffer.String(), errTruncated)

		} else if ctx.Err() == context.DeadlineExceeded {

			return "", sanitizeError(name, image, buffer.String(), errTimeout)

		} else if ctx.Err() == context.Canceled {

			return "", sanitizeError(name, image, buffer.String(), errIdle)

		} else if err1 == nil {

			return buffer.String(), nil

		} else {
			return "", sanitizeError(name, image, buffer.String(), err1)
		}

	} else {
		return "", sanitizeError(name, image, "", err0)
	}

}
