package markdown

import utils_strings "exocomp/utils/strings"
import "strings"

func generateId(element *Element) string {

	texts := make([]string, 0)

	if element.Text != "" {

		chunks := strings.Split(strings.TrimSpace(utils_strings.ToASCII(element.Text)), "-")

		for _, chunk := range chunks {

			tmp := strings.TrimSpace(chunk)

			if tmp != "" {
				texts = append(texts, tmp)
			}

		}

	} else if len(element.Children) > 0 {

		for _, child := range element.Children {

			if child.Type == "b" || child.Type == "code" || child.Type == "del" || child.Type == "em" || child.Type == "#text" {

				tmp := strings.TrimSpace(utils_strings.ToASCIIName(child.Text))

				if strings.HasPrefix(tmp, "-") {
					tmp = tmp[1:]
				}

				if strings.HasSuffix(tmp, "-") {
					tmp = tmp[0 : len(tmp)-1]
				}

				chunks := strings.Split(tmp, "-")

				for _, chunk := range chunks {

					tmp := strings.TrimSpace(chunk)

					if tmp != "" {
						texts = append(texts, tmp)
					}

				}

			}

		}

	}

	if len(texts) > 0 {

		filtered := make([]string, 0)

		if utils_strings.IsNumber(string(texts[0])) {
			texts = texts[1:]
		}

		for _, text := range texts {

			tmp := strings.ToLower(strings.TrimSpace(text))

			if tmp != "" {
				filtered = append(filtered, tmp)
			}

		}

		return strings.Join(filtered, "-")

	}

	return ""

}
