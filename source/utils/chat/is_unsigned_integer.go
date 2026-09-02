package chat

import "reflect"

func is_unsigned_integer(k reflect.Kind) bool {

	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	default:
		return false
	}

}

