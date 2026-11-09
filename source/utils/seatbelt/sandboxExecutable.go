package seatbelt

import "fmt"
import "os"
import "os/exec"

const sandboxExecPath = "/usr/bin/sandbox-exec"

func sandboxExecutable() (string, error) {

	stat, err0 := os.Stat(sandboxExecPath)

	if err0 == nil && stat.Mode().IsRegular() == true {

		return sandboxExecPath, nil

	} else {

		path, err1 := exec.LookPath("sandbox-exec")

		if err1 == nil {
			return path, nil
		} else {
			return "", fmt.Errorf("seatbelt.sandboxExecutable: sandbox-exec is not installed.")
		}

	}

}
