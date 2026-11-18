package utils

import "exocomp-toolchain/types"
import "encoding/json"
import "fmt"
import "os"
import "path/filepath"

func ReadPrograms(base_dir string) ([]types.Program, error) {

	path := filepath.Join(base_dir, "toolchain", "programs.json")
	buffer, err0 := os.ReadFile(path)

	if err0 == nil {

		programs := make([]types.Program, 0)
		err1 := json.Unmarshal(buffer, &programs)

		if err1 == nil {

			err2 := ValidatePrograms(programs)

			if err2 == nil {
				return programs, nil
			} else {
				return nil, err2
			}

		} else {
			return nil, fmt.Errorf("utils.ReadPrograms: Invalid manifest: %s", err1.Error())
		}

	} else {
		return nil, fmt.Errorf("utils.ReadPrograms: Cannot read manifest \"toolchain/programs.json\".")
	}

}
