package config

import "fmt"
import "strings"
import "exocomp/gadgets"

type Gadget struct {
	Type      GadgetType
	Method    GadgetMethod
	Arguments GadgetArguments
	Lines     [2]int
}

func ParseGadgets(text string) []*Gadget {

	gadgets := make([]*Gadget, 0)
	lines   := strings.Split(text, "\n")

	for l := 0; l < len(lines); l++ {

		line = strings.TrimSpace(lines[l])

		if strings.HasPrefix(line, "#!gadget:" + GadgetTypeHelp.String() + ".") {

			// #!gadget:help.Overview
			// #!gadget:help.Gadget <gadget>
			fields := strings.Fields(strings.TrimSpace(line[len(GadgetTypeHelp.String()) + 10:]))
			method := fields[0]

			if strings.ToLower(method) == "overview" {

				arguments := fields[1:]

				gadgets = append(gadgets, &Gadget{
					Type:      GadgetTypeHelp,
					Method:    GadgetMethod("Overview"),
					Arguments: GadgetArguments(arguments),
					Lines:     [2]int{l+1, l+1},
				})

			} else if strings.ToLower(method) == "gadget" {

				arguments := fields[1:]

				gadgets = append(gadgets, &Gadget{
					Type:      GadgetTypeHelp,
					Method:    GadgetMethod("Gadget"),
					Arguments: GadgetArguments(arguments),
					Lines:     [2]int{l+1, l+1},
				})

			} else {

				arguments := fields[1:]

				gadgets = append(gadgets, &Gadget{
					Type:      GadgetTypeInvalid,
					Method:    GadgetMethod(method),
					Arguments: GadgetArguments(arguments),
					Lines:     [2]int{l+1, l+1},
				})

			}

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypeFiles.String() + ".") {

			// TODO: Parse from here to the last line

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypePrograms.String() + ".") {

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypeTasks.String() + ".") {

			// TODO: Implement tasks.List()
			// TODO: tasks.Add()
			// TODO: tasks.Done()

		}

	}

	return gadgets

}

func (gadget *Gadget) Call(config *Config) (string, error) {

	// TODO: Implement delegation to actual APIs
	// Is it possible to use an Interface here?

	switch gadget.Type {
	case GadgetTypeHelp:

		help_gadget := gadgets.NewHelp(config.Sandbox, config.Gadgets)

		if gadget.Method == "Overview" {

			return help_gadget.Help(gadget.Arguments)

		} else if gadget.Method == "Gadget" {

			topic := strings.ToLower(gadget.Arguments.Get(0))

			if topic == GadgetTypeHelp.String() {
				return help_gadget.Help(gadget.Arguments)
			} else if topic == GadgetTypeFiles.String() {
				return gadgets.NewFiles(config.Sandbox).Help(gadget.Arguments)
			} else if topic == GadgetTypePrograms.String() {
				return gadgets.NewPrograms(config.Sandbox, config.Programs).Help(gadget.Arguments)
			} else if topic == GadgetTypeTasks.String() {
				return gadgets.NewTasks(config.Sandbox).Help(gadget.Arguments)
			}

		}

	case GadgetTypeFiles:

		files_gadget := gadgets.NewFiles(config.Sandbox)

		if gadget.Method == "List" {
			return files_gadget.List(gadget.Arguments)
		} else if gadget.Method == "Read" {
			return files_gadget.Read(gadget.Arguments)
		} else if gadget.Method == "Stat" {
			return files_gadget.Stat(gadget.Arguments)
		} else if gadget.Method == "Write" {
			return files_gadget.Write(gadget.Arguments)
		}

	case GadgetTypePrograms:

		programs_gadget := gadgets.NewPrograms(config.Sandbox, config.Programs)

		if gadget.Method == "List" {
			return programs_gadget.List(gadget.Arguments)
		} else if gadget.Method == "Execute" {
			return programs_gadget.Execute(gadget.Arguments)
		}

	case GadgetTypeRevisions:

		// TODO

	case GadgetTypeTasks:

		// TODO

	}

	return "", fmt.Errorf("Gadget %s.%s is not implemented yet", gadget.Type.String(), gadget.Method.String())

}
