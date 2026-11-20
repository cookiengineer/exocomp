package components

func footerBindings(status string) []footerBinding {

	if status != "" {
		return []footerBinding{
			{Key: "Ctrl+↑↓", Label: "agents"},
			{Key: "Alt+↑↓", Label: "scroll"},
			{Key: "Ctrl+Q", Label: "quit"},
		}
	} else {
		return []footerBinding{
			{Key: "Ctrl+↑↓", Label: "agents"},
			{Key: "Alt+↑↓", Label: "scroll"},
			{Key: "Ctrl+S", Label: "send"},
			{Key: "/", Label: "tools"},
			{Key: "Ctrl+Q", Label: "quit"},
		}
	}

}
