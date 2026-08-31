package yaml

import "fmt"
import net_url "net/url"
import "reflect"

var net_url_type = reflect.TypeOf(net_url.URL{})
var net_url_pointer_type = reflect.TypeOf((*net_url.URL)(nil))

func decodeURL(node *Node, target reflect.Value) error {

	if node.Kind != ScalarNode {
		return fmt.Errorf(
			"cannot decode %v into net.URL",
			node.Kind,
		)
	}

	parsed, err := net_url.Parse(node.Value)

	if err != nil {
		return err
	}

	switch target.Kind() {

	case reflect.Struct:
		target.Set(reflect.ValueOf(*parsed))

	case reflect.Ptr:
		target.Set(reflect.ValueOf(parsed))

	default:
		return fmt.Errorf(
			"cannot decode net.URL into %s",
			target.Kind(),
		)

	}

	return nil

}
