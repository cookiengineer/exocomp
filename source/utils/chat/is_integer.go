package chat

import "reflect"

func is_integer(k reflect.Kind) bool {
	return is_signed_integer(k) || is_unsigned_integer(k)
}

