package chat

import "reflect"

func isInteger(k reflect.Kind) bool {
	return isSignedInteger(k) || isUnsignedInteger(k)
}
