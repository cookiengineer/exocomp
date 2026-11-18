package actions

import "os"
import "strings"

func isProgramCached(cache_output string, version string) bool {

	stat, err0 := os.Stat(cache_output)

	if err0 == nil && stat.Mode().IsRegular() == true {

		buffer, err1 := os.ReadFile(cache_output + ".version")

		if err1 == nil {
			return strings.TrimSpace(string(buffer)) == version
		} else {
			return false
		}

	} else {
		return false
	}

}
