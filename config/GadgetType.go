package config

type GadgetType string

const (
	GadgetTypeHelp     GadgetType = "help"
	GadgetTypeFiles    GadgetType = "files"
	GadgetTypePrograms GadgetType = "programs"
	GadgetTypeTasks    GadgetType = "tasks"
)

func IsGadgetType(raw string) bool {

	switch GadgetType(raw) {
	case GadgetTypeHelp:
		return true
	case GadgetTypeFiles:
		return true
	case GadgetTypePrograms:
		return true
	case GadgetTypeTasks:
		return true
	}

	return false

}

func (gadget_type GadgetType) String() string {
	return string(gadget_type)
}
