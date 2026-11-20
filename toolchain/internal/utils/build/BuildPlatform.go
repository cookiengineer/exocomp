package build

import utils_assets "exocomp-toolchain/internal/utils/assets"
import "fmt"
import "os"

func BuildPlatform(base_dir string, operating_system string, architecture string, is_building_programs bool, is_building_installer bool) error {

	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "== Build %s ==\n", operating_system)

	utils_assets.ClearAssets(base_dir)
	BuildExocomps(base_dir, operating_system, architecture)

	if is_building_programs == true {

		err := BuildPrograms(base_dir, operating_system, architecture)

		if err != nil {
			return err
		}

	}

	if is_building_installer == true {
		BuildInstaller(base_dir, operating_system, architecture)
	}

	return nil

}
