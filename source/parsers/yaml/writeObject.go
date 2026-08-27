package yaml

import "fmt"
import "strings"

func writeObject(builder *strings.Builder, node *Node, indent int) {

	indentation := strings.Repeat(" ", indent)

	for key, child := range node.ObjectChildren {

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
