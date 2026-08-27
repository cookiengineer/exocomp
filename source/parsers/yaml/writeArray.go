package yaml

import "fmt"
import "strings"

func writeArray(builder *strings.Builder, node *Node, indent int) {

	indentation := strings.Repeat(" ", indent)

	for _, child := range node.ArrayChildren {

		builder.WriteString(
			fmt.Sprintf(
				"%s- %s\n",
				indentation,
				child.Value,
			),
		)

	}

}
