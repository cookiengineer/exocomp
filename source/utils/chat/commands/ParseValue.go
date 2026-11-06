package commands

import "exocomp/schemas"
import "encoding/json"
import "fmt"
import "strconv"
import "strings"

func ParseValue(raw string, property *schemas.ToolFunctionParameterProperty) (any, error) {

	value := strings.TrimSpace(raw)
	property_type := ""

	if property != nil {
		property_type = strings.ToLower(strings.TrimSpace(property.Type))
	}

	if property_type == "string" || property_type == "" {

		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {

			unquoted, err := strconv.Unquote(value)

			if err == nil {
				value = unquoted
			} else {
				value = value[1 : len(value)-1]
			}

		} else if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
			value = value[1 : len(value)-1]
		} else if property_type == "" {
			return parseSchemaLess(value)
		}

		if property != nil && len(property.Enum) > 0 {

			found := false

			for _, option := range property.Enum {

				if option == value {
					found = true
					break
				}

			}

			if found == true {
				return value, nil
			} else {
				return nil, fmt.Errorf("commands.ParseValue: \"%s\" is not one of %s.", value, strings.Join(property.Enum, ", "))
			}

		} else {
			return value, nil
		}

	}

	if property_type == "integer" {

		parsed, err := strconv.ParseInt(unquote(value), 10, 64)

		if err == nil {
			return parsed, nil
		} else {
			return nil, fmt.Errorf("commands.ParseValue: \"%s\" is not an integer.", value)
		}

	}

	if property_type == "number" {

		parsed, err := strconv.ParseFloat(unquote(value), 64)

		if err == nil {
			return parsed, nil
		} else {
			return nil, fmt.Errorf("commands.ParseValue: \"%s\" is not a number.", value)
		}

	}

	if property_type == "boolean" {

		parsed, err := strconv.ParseBool(unquote(value))

		if err == nil {
			return parsed, nil
		} else {
			return nil, fmt.Errorf("commands.ParseValue: \"%s\" is not a boolean.", value)
		}

	}

	if property_type == "null" {
		return nil, nil
	}

	if property_type == "array" || property_type == "object" || strings.HasPrefix(value, "[") || strings.HasPrefix(value, "{") {

		var decoded any
		err := json.Unmarshal([]byte(value), &decoded)

		if err == nil {
			return decoded, nil
		} else {
			return nil, fmt.Errorf("commands.ParseValue: \"%s\" is not valid JSON.", value)
		}

	}

	return parseSchemaLess(value)

}

func parseSchemaLess(value string) (any, error) {

	if value == "" {
		return "", nil
	}

	if value == "true" {
		return true, nil
	}

	if value == "false" {
		return false, nil
	}

	if value == "null" {
		return nil, nil
	}

	var decoded any
	err := json.Unmarshal([]byte(value), &decoded)

	if err == nil {
		return decoded, nil
	} else {
		return unquote(value), nil
	}

}

func unquote(value string) string {

	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {

		unquoted, err := strconv.Unquote(value)

		if err == nil {
			return unquoted
		}

	}

	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}

	return value

}
