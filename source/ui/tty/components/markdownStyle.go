package components

import utils_tty "exocomp/utils/tty"

func markdownStyle(styles map[string]utils_tty.Style, kind string) utils_tty.Style {

	style, ok := styles[kind]

	if ok == true {
		return style
	} else {
		return styles["text"]
	}

}
