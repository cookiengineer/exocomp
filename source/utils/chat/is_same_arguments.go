package chat

import "reflect"

func is_same_arguments(a map[string]any, b map[string]any) bool {

	if len(a) != len(b) {
		return false
	}

	for key, av := range a {

		bv, ok := b[key]
		if !ok {
			return false
		}

		if !is_same_value(reflect.ValueOf(av), reflect.ValueOf(bv)) {
			return false
		}

	}

	return true

}

