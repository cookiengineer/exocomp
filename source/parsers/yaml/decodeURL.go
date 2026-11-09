package yaml

import "fmt"
import net_url "net/url"
import "reflect"

var netURLType = reflect.TypeOf(net_url.URL{})
var netURLPointerType = reflect.TypeOf((*net_url.URL)(nil))

func decodeURL(node *Node, target reflect.Value) error {

	if node.Kind != ScalarNode {
		return fmt.Errorf(
			"yaml.DecodeURL: Cannot decode %v into net.URL.",
			node.Kind,
		)
	}

	parsed, err := net_url.Parse(node.Value)

	if err == nil {

		switch target.Kind() {

		case reflect.Struct:
			target.Set(reflect.ValueOf(*parsed))

		case reflect.Ptr:
			target.Set(reflect.ValueOf(parsed))

		default:
			return fmt.Errorf(
				"yaml.DecodeURL: Cannot decode net.URL into %s.",
				target.Kind(),
			)

		}

	} else {
		return err
	}

	return nil

}
