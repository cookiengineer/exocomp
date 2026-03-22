package ollama

import "errors"
import "strings"
import "exocomp/gadgets"
import "exocomp/types"

type Gadget struct {
	Type      GadgetType
	Method    GadgetMethod
	Arguments []string
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

		if strings.HasPrefix(line, "#!" + GadgetTypeHelp.String() + " ") {

			// #!go-go-gadgeto-help <topic> <subtopic> ...
			arguments := strings.Fields(strings.TrimSpace(line[len(GadgetTypeHelp.String()) + 3:]))

			return &Gadget{
				Type:      GadgetTypeHelp,
				Method:    GadgetMethod("Help"),
				Arguments: GadgetArguments(arguments),
			}

		} else if strings.HasPrefix(line, "#!" + GadgetTypePrograms.String() + " ") {

			// #!go-go-gadgeto-program <program> <arg1> <arg2> <arg3> ...
			arguments := strings.Fields(strings.TrimSpace(line[len(GadgetTypePrograms.String()) + 3:]))

			return &Gadget{
				Type:      GadgetTypePrograms,
				Method:    GadgetMethod("Execute"),
				Arguments: GadgetArguments(arguments),
			}

		} else if strings.HasPrefix(line, "#!" + GadgetTypePrograms.String() + " ") {

			// #!go-go-gadgeto-filesystem <method> <file>
			fields := strings.Fields(strings.TrimSpace(line[len(GadgetTypeFiles.String()) + 3:]))
			method := fields[0]

			for f := 1; f < len(fields); f++ {

				if fields[f] == "<<" && f < len(fields) - 1 {

					seek := fields[f + 1]

					for s := l + 1; s < len(lines); s++ {

						if strings.HasPrefix(lines[s], seek) {

							fields[f] = strings.Join(lines[l+1:s], "\n")
							fields = append(fields[:f+1], fields[f+2:]...)
							f--
							break

						}

					}

				}

			}

			if method == "Read" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("Read"),
					Arguments: GadgetArguments(arguments),
				}

			} else if method == "Stat" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("Stat"),
					Arguments: GadgetArguments(arguments),
				}

			} else if method == "Write" {

				arguments := fields[1:]

				return &Gadget{
					Type:      GadgetTypeFiles,
					Method:    GadgetMethod("Write"),
					Arguments: GadgetArguments(arguments),
				}

			}

		} else if strings.HasPrefix(line, "#!" + GadgetTypeRevisions.String() + " ") {

			// TODO: commit, add, remove

		} else if strings.HasPrefix(line, "#!" + GadgetTypeTasks.String() + " ") {

			// TODO: Implement tasks.List()
			// TODO: tasks.Add()
			// TODO: tasks.Done()

		}

	}

	return nil

}

func (gadget *Gadget) Execute(config *types.Config) (string, error) {

	// TODO: Implement delegation to actual APIs
	// Is it possible to use an Interface here?

	switch gadget.Type {
	case GadgetTypeHelp:

		help_gadget := gadgets.NewHelp(config)

		if gadget.Method == "Help" {
			return help_gadget.Help(gadget.Arguments)
		}

	case GadgetTypeFiles:

		files_gadget := gadgets.NewFiles(config)

		if gadget.Method == "Read" {
			return files_gadget.Read(gadget.Arguments)
		} else if gadget.Method == "Stat" {
			return files_gadget.Stat(gadget.Arguments)
		} else if gadget.Method == "Write" {
			return files_gadget.Write(gadget.Arguments)
		}

	case GadgetTypePrograms:

		// TODO

	case GadgetTypeRevisions:

		// TODO

	case GadgetTypeTasks:

		// TODO

	}

	return "", fmt.Errorf("Gadget %s.%s is not implemented yet", gadget.Type.String(), gadget.Method.String())

}
