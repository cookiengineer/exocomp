package chat

import "reflect"

func isSameArguments(a map[string]any, b map[string]any) bool {

	if len(a) != len(b) {
		return false
	}

	for key, av := range a {

		bv, ok := b[key]
		if ok != true {
			return false
		}

		if isSameValue(reflect.ValueOf(av), reflect.ValueOf(bv)) == false {
			return false
		}

	}

	return true

}
