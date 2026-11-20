package lint

import "strings"

func isCamelLocal(name string) bool {

	if camelCasePattern.MatchString(name) == false {
		return false
	} else {

		segments := camelSegments(name)

		for _, segment := range segments[1:] {

			is_upper := segment == strings.ToUpper(segment)

			if is_upper == true && allowedLocalAcronyms[segment] == true {
				continue
			}

			return true

		}

		return false

	}

}
