package lint

func checkTypeInitialism(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, type_spec := range declaredTypes(file) {

		name := type_spec.Name.Name

		if initialismPattern.MatchString(name) == false {
			continue
		}

		if allowedInitialisms[name] == true {
			continue
		}

		result = append(
			result,
			file.Diagnostic(
				"naming/type-initialism",
				type_spec.Pos(),
				"type name %s is collapsed initials; spell it out",
				name,
			),
		)

	}

	return result

}
