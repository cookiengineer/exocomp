package yaml

import "fmt"
import "reflect"

func decodeMap(node *Node, target reflect.Value) error {

	if node.Kind != ObjectNode {

		return fmt.Errorf(
			"cannot decode %v into map",
			node.Kind,
		)

	}

	if target.Kind() != reflect.Map {

		return fmt.Errorf(
			"target is not a map: %s",
			target.Kind(),
		)

	}

	if target.Type().Key().Kind() != reflect.String {

		return fmt.Errorf(
			"unsupported map key type: %s",
			target.Type().Key().Kind(),
		)

	}

	if target.IsNil() == true {

		target.Set(
			reflect.MakeMap(target.Type()),
		)

	}

	map_value_type := target.Type().Elem()

	for key, child := range node.ObjectChildren {

		map_value := reflect.New(map_value_type).Elem()

		err := decodeValue(child, map_value)

		if err != nil {
			return err
		}

		target.SetMapIndex(
			reflect.ValueOf(key),
			map_value,
		)

	}

	return nil

}
