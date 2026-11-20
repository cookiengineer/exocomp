package chat

// Arguments is the decoded JSON argument map of a tool call. The typed
// accessors keep tool implementations from repeating type assertions.
type Arguments map[string]any

// Bool returns the boolean value of key.
func (arguments Arguments) Bool(key string) (bool, bool) {

	value, ok := arguments[key].(bool)

	if ok == true {
		return value, true
	} else {
		return false, false
	}

}

// Float returns the numeric value of key as a float64.
func (arguments Arguments) Float(key string) (float64, bool) {

	value, ok := arguments[key].(float64)

	if ok == true {
		return value, true
	} else {
		return 0, false
	}

}

// Has reports whether key is present.
func (arguments Arguments) Has(key string) bool {

	_, ok := arguments[key]

	return ok

}

// Int returns the integer value of key.
func (arguments Arguments) Int(key string) (int, bool) {

	value, ok := arguments[key]

	if ok == true {

		switch typed := value.(type) {
		case float64:
			return int(typed), true
		case int:
			return typed, true
		default:
			return 0, false
		}

	} else {
		return 0, false
	}

}

// Map returns the string map value of key.
func (arguments Arguments) Map(key string) (map[string]string, bool) {

	raw, ok := arguments[key]

	if ok == true {

		switch typed := raw.(type) {
		case map[string]string:
			return typed, true
		case map[string]any:

			result := make(map[string]string, len(typed))

			for name, value := range typed {

				if text, ok := value.(string); ok == true {
					result[name] = text
				}

			}

			return result, true

		default:
			return nil, false
		}

	} else {
		return nil, false
	}

}

// Object returns the nested object value of key.
func (arguments Arguments) Object(key string) (Arguments, bool) {

	raw, ok := arguments[key].(map[string]any)

	if ok == true {
		return Arguments(raw), true
	} else {
		return nil, false
	}

}

// String returns the string value of key.
func (arguments Arguments) String(key string) (string, bool) {

	value, ok := arguments[key].(string)

	if ok == true {
		return value, true
	} else {
		return "", false
	}

}

// Strings returns the string slice value of key.
func (arguments Arguments) Strings(key string) ([]string, bool) {

	raw, ok := arguments[key]

	if ok == true {

		switch typed := raw.(type) {
		case []string:
			return typed, true
		case []any:

			result := make([]string, 0, len(typed))

			for _, value := range typed {

				if text, ok := value.(string); ok == true {
					result = append(result, text)
				}

			}

			return result, true

		default:
			return nil, false
		}

	} else {
		return nil, false
	}

}
