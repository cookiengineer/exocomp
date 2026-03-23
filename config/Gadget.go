package config

import "fmt"
import "strings"
import "exocomp/gadgets"

type Gadget struct {
	Type      GadgetType
	Method    GadgetMethod
	Arguments GadgetArguments
}

func UsesGadget(text string) bool {

	gadget := ParseGadget(text)

	if gadget != nil {
		return true
	}

	return false

}

func ParseGadget(text string) *Gadget {

	lines := strings.Split(text, "\n")

	for l, line := range lines {

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#!gadget:" + GadgetTypeHelp.String() + ".") {

			// #!gadget:help.Overview
			// #!gadget:help.Gadget <gadget>
			fields := strings.Fields(strings.TrimSpace(line[len(GadgetTypeHelp.String()) + 10:]))
			method := fields[0]

			if strings.ToLower(method) == "overview" {

				return &Gadget{
					Type:      GadgetTypeHelp,
					Method:    GadgetMethod("Overview"),
					Arguments: GadgetArguments([]string{}),
				}

			} else if strings.ToLower(method) == "gadget" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeHelp,
					Method:    GadgetMethod("Gadget"),
					Arguments: GadgetArguments(arguments),
				}

			}

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypeFiles.String() + ".") {

			// #!gadget:files.Read <path>
			// #!gadget:files.Stat <path>
			// #!gadget:files.Write <path> <<#!EOF
			// ...
			// #!EOF
			fields := strings.Fields(strings.TrimSpace(line[len(GadgetTypeFiles.String()) + 10:]))
			method := fields[0]

			for f, field := range fields {

				if strings.HasPrefix(field, "\"") && strings.HasSuffix(field, "\"") {
					fields[f] = field[1:len(field)-1]
				} else if strings.HasPrefix(field, "'") && strings.HasSuffix(field, "'") {
					fields[f] = field[1:len(field)-1]
				}

			}

			for f := 1; f < len(fields); f++ {

				if strings.HasPrefix(fields[f], "<<") {

					seek := fields[f][2:]

					for s := l + 1; s < len(lines); s++ {

						if strings.HasPrefix(lines[s], seek) {

							fields[f] = strings.Join(lines[l+1:s], "\n")
							break

						}

					}

				}

			}

			if strings.ToLower(method) == "list" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("List"),
					Arguments: GadgetArguments(arguments),
				}

			} else if strings.ToLower(method) == "read" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("Read"),
					Arguments: GadgetArguments(arguments),
				}

			} else if strings.ToLower(method) == "stat" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("Stat"),
					Arguments: GadgetArguments(arguments),
				}

			} else if strings.ToLower(method) == "write" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("Write"),
					Arguments: GadgetArguments(arguments),
				}

			} else {
				return nil
			}

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypePrograms.String() + ".") {

			// #!programs.Execute <arguments...>
			fields := strings.Fields(strings.TrimSpace(line[len(GadgetTypePrograms.String()) + 10:]))
			method := fields[0]

			if strings.ToLower(method) == "execute" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypePrograms,
					Method:    GadgetMethod("Execute"),
					Arguments: GadgetArguments(arguments),
				}

			} else {
				return nil
			}

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypeRevisions.String() + ".") {

			// TODO: commit, add, remove

		} else if strings.HasPrefix(line, "#!gadget:" + GadgetTypeTasks.String() + ".") {

			// TODO: Implement tasks.List()
			// TODO: tasks.Add()
			// TODO: tasks.Done()

		}

	}

	return nil

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
			} else if topic == GadgetTypeRevisions.String() {
				// TODO: return gadgets.NewRevisions(config.Sandbox).Help(gadget.Arguments)
			} else if topic == GadgetTypeTasks.String() {
				// TODO: return gadgets.NewTasks(config.Sandbox).Help(gadget.Arguments)
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
