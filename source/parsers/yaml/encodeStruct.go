package yaml

import "reflect"
import "strings"

func encodeStruct(value reflect.Value) (*Node, error) {

	node := &Node{
		Kind:           ObjectNode,
		ObjectChildren: map[string]*Node{},
	}

	value_type := value.Type()

	for index := 0; index < value.NumField(); index++ {

		field_value := value.Field(index)
		field_type  := value_type.Field(index)

		yaml_tag := parseYAMLTag(field_type.Tag.Get("yaml"))

		if yaml_tag == "" {
			yaml_tag = strings.ToLower(field_type.Name)
		}

		child_node, err := encodeValue(field_value)

		if err != nil {
			return nil, err
		}

		node.ObjectChildren[yaml_tag] = child_node

	}

	return node, nil

}
