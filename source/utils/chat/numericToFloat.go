package chat

import "reflect"

func numericToFloat(v reflect.Value) float64 {

	if isSignedInteger(v.Kind()) {
		return float64(v.Int())
	}

	return float64(v.Uint())

}
