package chat

import "reflect"

func is_float(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

