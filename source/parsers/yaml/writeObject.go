package yaml

import "fmt"
import "sort"
import "strings"

func writeObject(builder *strings.Builder, node *Node, indent int) {

	indentation := strings.Repeat(" ", indent)
	keys := make([]string, 0, len(node.ObjectChildren))

	for key := range node.ObjectChildren {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {

		child := node.ObjectChildren[key]

		switch child.Kind {

		case ScalarNode:

			builder.WriteString(
				fmt.Sprintf(
					"%s%s: %s\n",
					indentation,
					key,
					child.Value,
				),
			)

		case ObjectNode, ArrayNode:

			builder.WriteString(
				fmt.Sprintf(
					"%s%s:\n",
					indentation,
					key,
				),
			)

			writeNode(
				builder,
				child,
				indent+2,
			)

		}

	}

}
