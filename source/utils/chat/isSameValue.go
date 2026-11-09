package chat

import "reflect"

func isSameValue(a reflect.Value, b reflect.Value) bool {

	if a.IsValid() == false || b.IsValid() == false {
		return a.IsValid() == false && b.IsValid() == false
	}

	// Treat all numeric types as comparable by value.
	if isInteger(a.Kind()) && isInteger(b.Kind()) {

		if isSignedInteger(a.Kind()) && isSignedInteger(b.Kind()) {
			return a.Int() == b.Int()
		}

		if isUnsignedInteger(a.Kind()) && isUnsignedInteger(b.Kind()) {
			return a.Uint() == b.Uint()
		}

		if isSignedInteger(a.Kind()) {

			i := a.Int()
			return i >= 0 && uint64(i) == b.Uint()

		}

		u := a.Uint()

		return u <= (^uint64(0)>>1) && int64(u) == b.Int()

	}

	if isFloat(a.Kind()) && isFloat(b.Kind()) {
		return a.Float() == b.Float()
	}

	// Allow integers/floats to compare by numeric value.
	if isInteger(a.Kind()) && isFloat(b.Kind()) {
		return numericToFloat(a) == b.Float()
	}

	if isFloat(a.Kind()) && isInteger(b.Kind()) {
		return a.Float() == numericToFloat(b)
	}

	if a.Kind() != b.Kind() {
		return false
	}

	switch a.Kind() {

	case reflect.Bool, reflect.String:
		return a.Interface() == b.Interface()

	case reflect.Interface:
		return isSameValue(a.Elem(), b.Elem())

	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}

		for i := 0; i < a.Len(); i++ {
			if isSameValue(a.Index(i), b.Index(i)) == false {
				return false
			}
		}

		return true

	case reflect.Map:
		if a.Len() != b.Len() {
			return false
		}

		for _, key := range a.MapKeys() {

			av := a.MapIndex(key)
			bv := b.MapIndex(key)

			if bv.IsValid() == false || isSameValue(av, bv) == false {
				return false
			}

		}

		return true

	case reflect.Complex64, reflect.Complex128:
		return a.Complex() == b.Complex()

	case reflect.Pointer:
		return isSameValue(a.Elem(), b.Elem())

	default:

		// Covers other comparable native values.
		if a.Type() != b.Type() {
			return false
		}

		return reflect.DeepEqual(a.Interface(), b.Interface())

	}

}
