package lint

import "strings"

func isPredicateName(name string) bool {

	if predicateExceptions[name] == true {
		return false
	} else {

		for _, prefix := range predicatePrefixes {

			if strings.HasPrefix(name, prefix) == true {
				return false
			}

		}

		return true

	}

}
