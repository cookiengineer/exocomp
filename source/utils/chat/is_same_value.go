package chat

import "reflect"

func is_same_value(a reflect.Value, b reflect.Value) bool {

	if !a.IsValid() || !b.IsValid() {
		return !a.IsValid() && !b.IsValid()
	}

	// Treat all numeric types as comparable by value.
	if is_integer(a.Kind()) && is_integer(b.Kind()) {

		if is_signed_integer(a.Kind()) && is_signed_integer(b.Kind()) {
			return a.Int() == b.Int()
		}

		if is_unsigned_integer(a.Kind()) && is_unsigned_integer(b.Kind()) {
			return a.Uint() == b.Uint()
		}

		if is_signed_integer(a.Kind()) {
			i := a.Int()
			return i >= 0 && uint64(i) == b.Uint()
		}

		u := a.Uint()

		return u <= (^uint64(0)>>1) && int64(u) == b.Int()

	}

	if is_float(a.Kind()) && is_float(b.Kind()) {
		return a.Float() == b.Float()
	}

	// Allow integers/floats to compare by numeric value.
	if is_integer(a.Kind()) && is_float(b.Kind()) {
		return numeric_to_float(a) == b.Float()
	}

	if is_float(a.Kind()) && is_integer(b.Kind()) {
		return a.Float() == numeric_to_float(b)
	}

	if a.Kind() != b.Kind() {
		return false
	}

	switch a.Kind() {
	case reflect.Bool, reflect.String:
		return a.Interface() == b.Interface()

	case reflect.Interface:
		return is_same_value(a.Elem(), b.Elem())

	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}

		for i := 0; i < a.Len(); i++ {
			if !is_same_value(a.Index(i), b.Index(i)) {
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

			if !bv.IsValid() || !is_same_value(av, bv) {
				return false
			}
		}

		return true

	case reflect.Complex64, reflect.Complex128:
		return a.Complex() == b.Complex()

	case reflect.Pointer:
		return is_same_value(a.Elem(), b.Elem())

	default:

		// Covers other comparable native values.
		if a.Type() != b.Type() {
			return false
		}

		return reflect.DeepEqual(a.Interface(), b.Interface())

	}

}

