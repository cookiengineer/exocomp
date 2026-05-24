package agents

import "exocomp/types"
import "embed"
import "io/fs"
import "path/filepath"

var Roles map[string]*types.Agent

//go:embed *.yaml
var filesystem embed.FS

func init() {

	Roles = make(map[string]*types.Agent)

	entries, err0 := fs.ReadDir(filesystem, ".")

	if err0 == nil {

		for _, entry := range entries {

			name := entry.Name()
			ext  := filepath.Ext(name)

			if ext == ".yaml" {

				data, err1 := filesystem.ReadFile(name)

				if err1 == nil {

					agent, err := types.ParseAgent(data)

					if err == nil {

						if agent.Role != "" {
							Roles[agent.Role] = agent
						}

					}

				}

			}

		}

	}

}

