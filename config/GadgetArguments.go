package config

type GadgetArguments []string

func (gadget_arguments GadgetArguments) Get(a int) string {

	if a >= 0 && a < len(gadget_arguments) {
		return string(gadget_arguments[a])
	}

	return ""

}

func (gadget_arguments GadgetArguments) Strings() []string {
	return []string(gadget_arguments)
}

