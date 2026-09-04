package engine

import "exocomp/engine"
import "exocomp/types"
import "os"
import "path/filepath"

func RestoreSession(source string, config *types.Config) (*engine.Session, error) {

	bytes, err1 := os.ReadFile(filepath.Join(source, ".exocomp", "session.json"))

	if err1 == nil {

		backup, err2 := engine.ParseSession(bytes)

		if err2 == nil {

			if config != nil {

				backup.Config.Name       = config.Name
				backup.Config.Role       = config.Role
				backup.Config.Playground = config.Playground
				backup.Config.Sandbox    = config.Sandbox

			}

			return engine.RestoreSession(backup.Config.Playground, *backup), nil

		} else {
			return nil, err2
		}

	} else {
		return nil, err1
	}

}

