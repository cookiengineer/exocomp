package gadgets

import "strings"

type Help struct {
	Gadgets []string
	Sandbox string
}

func NewHelp(sandbox string, gadgets []string) *Help {

	return &Help{
		Gadgets: gadgets,
		Sandbox: sandbox,
	}

}

func (gadget *Help) Help(arguments []string) (string, error) {

	gadget_help := make([]string, 0)

	for _, name := range gadget.Gadgets {
		gadget_help = append(gadget_help, "#!gadget:help.Gadget " + name)
	}

	return strings.Join([]string{
		"Show this overview:",
		"#!gadget:help.Overview",
		"",
		"Show detailed gadget help:",
		strings.Join(gadget_help, "\n"),
	}, "\n"), nil

}
