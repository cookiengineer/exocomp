package utils

import "exocomp-toolchain/types"
import "fmt"
import "strings"

func ValidatePrograms(programs []types.Program) error {

	if len(programs) == 0 {
		return fmt.Errorf("utils.ValidatePrograms: Manifest is empty.")
	}

	names := make(map[string]bool, len(programs))

	for _, program := range programs {

		if program.Name == "" || strings.ContainsAny(program.Name, "/\\") == true {
			return fmt.Errorf("utils.ValidatePrograms: Invalid program name \"%s\".", program.Name)
		} else if program.Folder == "" || strings.ContainsAny(program.Folder, "/\\") == true {
			return fmt.Errorf("utils.ValidatePrograms: Invalid folder for \"%s\".", program.Name)
		} else if program.Repository == "" {
			return fmt.Errorf("utils.ValidatePrograms: Missing repository for \"%s\".", program.Name)
		} else if program.Version == "" {
			return fmt.Errorf("utils.ValidatePrograms: Missing version for \"%s\".", program.Name)
		} else if program.Source == "" {
			return fmt.Errorf("utils.ValidatePrograms: Missing source for \"%s\".", program.Name)
		}

		if names[program.Name] == true {
			return fmt.Errorf("utils.ValidatePrograms: Duplicate program \"%s\".", program.Name)
		} else {
			names[program.Name] = true
		}

	}

	return nil

}
