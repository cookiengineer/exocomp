package main

import "installer/actions"
import "os"
import "strings"

func main() {

	// prefix := "/opt/exocomp"
	prefix := "/usr/local"

	if len(os.Args) > 1 {

		if strings.HasPrefix(os.Args[1], "--prefix=") {

			tmp := strings.TrimSpace(os.Args[1][9:])

			if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
				tmp = strings.TrimSpace(tmp[1:len(tmp)-1])
			}

			if strings.HasPrefix(tmp, "'") && strings.HasSuffix(tmp, "'") {
				tmp = strings.TrimSpace(tmp[1:len(tmp)-1])
			}

			if strings.HasPrefix(tmp, "/") && len(tmp) > 1 {
				prefix = tmp
			}

		}

	}

	actions.Install(prefix)

}
