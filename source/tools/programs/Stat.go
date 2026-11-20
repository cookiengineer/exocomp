package programs

import "fmt"
import "os/exec"
import "path/filepath"

func (tool *Programs) Stat(program string) (string, error) {

	found := false

	for _, name := range tool.AllowedPrograms {

		if name == program {

			found = true
			break

		}

	}

	if found == true {

		path, err1 := exec.LookPath(program)

		if err1 == nil {

			real_path, err2 := filepath.EvalSymlinks(path)

			if err2 == nil {
				return real_path, nil
			} else {
				return "", fmt.Errorf("programs.Stat: Invalid program \"%s\": Program doesn't exist.", program)
			}

		} else {
			return "", fmt.Errorf("programs.Stat: Invalid program \"%s\": Program doesn't exist.", program)
		}

	} else {
		return "", fmt.Errorf("programs.Stat: Invalid program \"%s\": Attempt to lookup unallowed program", program)
	}

}
