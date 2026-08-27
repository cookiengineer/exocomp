package yaml

import "reflect"
import "strings"

func decodeStruct(node *Node, target reflect.Value) error {

	target_type := target.Type()

	for index := 0; index < target.NumField(); index++ {

		struct_field := target_type.Field(index)
		target_field := target.Field(index)

		yaml_tag := parseYAMLTag(struct_field.Tag.Get("yaml"))

		if yaml_tag == "" {
			yaml_tag = strings.ToLower(struct_field.Name)
		}

		if yaml_tag != "-" {

			child, exists := node.ObjectChildren[yaml_tag]

			if exists == true {

				err := decodeValue(child, target_field)

				if err != nil {
					return err
				}

			} else {
				continue
			}

		} else {
			continue
		}

	}

	return nil

}
