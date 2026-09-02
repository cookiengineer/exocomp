package chat

import "reflect"

func numeric_to_float(v reflect.Value) float64 {

	if is_signed_integer(v.Kind()) {
		return float64(v.Int())
	}

	return float64(v.Uint())

}

