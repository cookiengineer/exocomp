package build

import "fmt"
import "os"

func BuildTargets(base_dir string, operating_systems []string, architecture string, is_building_programs bool, is_building_installer bool) error {

	if is_building_programs == true {

		fmt.Fprint(os.Stdout, "\n")
		fmt.Fprint(os.Stdout, "== Clone programs ==\n")

		err := ClonePrograms(base_dir)

		if err != nil {
			return err
		}

	}

	for _, operating_system := range operating_systems {

		err := BuildPlatform(base_dir, operating_system, architecture, is_building_programs, is_building_installer)

		if err != nil {
			return err
		}

	}

	return nil

}
