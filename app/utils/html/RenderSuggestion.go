package html

import "fmt"
import "strings"

func RenderSuggestion(suggestion map[string]any) string {

	label, _ := suggestion["label"].(string)
	description, _ := suggestion["description"].(string)

	if strings.HasPrefix(label, "/") == false {
		label = "&nbsp;&nbsp;" + label
	}

	return fmt.Sprintf("<li><label>%s</label> <span>%s</span></li>", label, description)

}
