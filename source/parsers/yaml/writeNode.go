package yaml

import "strings"

func writeNode(builder *strings.Builder, node *Node, indent int) {

	switch node.Kind {

	case ObjectNode:
		writeObject(builder, node, indent)

	case ArrayNode:
		writeArray(builder, node, indent)

	case ScalarNode:
		builder.WriteString(node.Value)

	}

}
