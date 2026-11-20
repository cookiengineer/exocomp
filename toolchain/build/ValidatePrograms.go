package main

import "fmt"
import "strings"

func ValidatePrograms(programs []Program) error {

	if len(programs) > 0 {

		names := make(map[string]bool, len(programs))

		for _, program := range programs {

			if program.Name == "" || strings.ContainsAny(program.Name, "/\\") == true {
				return fmt.Errorf("build.ValidatePrograms: Invalid program name \"%s\".", program.Name)
			} else if program.Folder == "" || strings.ContainsAny(program.Folder, "/\\") == true {
				return fmt.Errorf("build.ValidatePrograms: Invalid folder for \"%s\".", program.Name)
			} else if program.Repository == "" {
				return fmt.Errorf("build.ValidatePrograms: Missing repository for \"%s\".", program.Name)
			} else if program.Version == "" {
				return fmt.Errorf("build.ValidatePrograms: Missing version for \"%s\".", program.Name)
			} else if program.Source == "" {
				return fmt.Errorf("build.ValidatePrograms: Missing source for \"%s\".", program.Name)
			}

			if names[program.Name] == true {
				return fmt.Errorf("build.ValidatePrograms: Duplicate program \"%s\".", program.Name)
			} else {
				names[program.Name] = true
			}

		}

		return nil

	} else {
		return fmt.Errorf("build.ValidatePrograms: Manifest is empty.")
	}

}
