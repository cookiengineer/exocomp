package chat

import "reflect"

func is_signed_integer(k reflect.Kind) bool {

	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}

}

