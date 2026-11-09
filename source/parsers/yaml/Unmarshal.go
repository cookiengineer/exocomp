package yaml

import "errors"
import "reflect"

func Unmarshal(data []byte, target any) error {

	parser := NewParser(data)
	root, err := parser.Root()

	if err == nil {

		target_value := reflect.ValueOf(target)

		if target_value.Kind() == reflect.Pointer {

			return decodeValue(
				root,
				target_value.Elem(),
			)

		} else {

			return errors.New(
				"yaml.Unmarshal: Target must be a pointer.",
			)

		}

	} else {
		return err
	}

}
