package seatbelt

import "fmt"
import "strings"

const maxProfileSize = 65535
const resolverSocket = "/private/var/run/mDNSResponder"

type ProfileBuilder struct {
	imports     []string
	rules       []string
	seenImports map[string]bool
}

func NewProfileBuilder() *ProfileBuilder {
	return &ProfileBuilder{
		imports:     make([]string, 0),
		rules:       make([]string, 0),
		seenImports: make(map[string]bool, 0),
	}
}

func (builder *ProfileBuilder) addImport(profile string) {

	if builder.seenImports[profile] == false {

		builder.seenImports[profile] = true
		builder.imports = append(builder.imports, fmt.Sprintf("(import %s)", quoteSBPL(profile)))

	}

}

func (builder *ProfileBuilder) addRule(rule string) {
	builder.rules = append(builder.rules, rule)
}

func (builder *ProfileBuilder) addDenyRules(paths []string, operations string) error {

	for _, path := range paths {

		err0 := validatePath(path)

		if err0 != nil {
			return err0
		}

	}

	resolved := resolvePaths(paths)

	if len(resolved) == 0 {
		return nil
	} else {

		filters := pathFilters(resolved)

		builder.addRule(fmt.Sprintf("(deny %s %s)", operations, filters))

		return nil

	}

}

func (builder *ProfileBuilder) addFileRules(paths []string, operations string, extension_classes []string) error {

	for _, path := range paths {

		err0 := validatePath(path)

		if err0 != nil {
			return err0
		}

	}

	resolved := resolvePaths(paths)

	if len(resolved) == 0 {
		return nil
	} else {

		filters := pathFilters(resolved)

		builder.addRule(fmt.Sprintf("(allow %s %s)", operations, filters))

		if len(extension_classes) > 0 {

			extension_filters := filters

			if len(resolved) > 1 {
				extension_filters = "(require-any " + filters + ")"
			}

			for _, extension_class := range extension_classes {

				builder.addRule(fmt.Sprintf(
					"(allow file-issue-extension (require-all (extension-class %s) %s))",
					quoteSBPL(extension_class),
					extension_filters,
				))

			}

		}

		return nil

	}

}

func (builder *ProfileBuilder) addProcessExecRules(paths []string) error {

	for _, path := range paths {

		err0 := validatePath(path)

		if err0 != nil {
			return err0
		}

	}

	resolved := resolvePaths(paths)

	if len(resolved) == 0 {
		return nil
	} else {

		builder.addRule(fmt.Sprintf("(allow process-exec %s)", pathFilters(resolved)))

		return nil

	}

}

func (builder *ProfileBuilder) build() (*Profile, error) {

	source_builder := strings.Builder{}

	source_builder.WriteString("(version 1)\n")
	source_builder.WriteString("(deny default)\n")

	for _, imported := range builder.imports {

		source_builder.WriteString(imported)
		source_builder.WriteString("\n")

	}

	for _, rule := range builder.rules {

		source_builder.WriteString(rule)
		source_builder.WriteString("\n")

	}

	source := source_builder.String()

	if len(source) > maxProfileSize {
		return nil, fmt.Errorf("seatbelt.BuildProfile: Profile exceeds 65535 byte limit.")
	} else {
		return &Profile{source: source}, nil
	}

}
