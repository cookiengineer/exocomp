package yaml

import "fmt"
import "reflect"

func encodeMap(value reflect.Value) (*Node, error) {

	node := &Node{
		Kind:           ObjectNode,
		ObjectChildren: map[string]*Node{},
	}

	for _, map_key := range value.MapKeys() {

		if map_key.Kind() != reflect.String {

			return nil, fmt.Errorf(
				"unsupported map key type: %s",
				map_key.Kind(),
			)

		}

		map_value       := value.MapIndex(map_key)
		child_node, err := encodeValue(map_value)

		if err != nil {
			return nil, err
		}

		node.ObjectChildren[map_key.String()] = child_node

	}

	return node, nil

}
