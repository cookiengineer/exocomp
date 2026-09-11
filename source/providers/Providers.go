package providers

import "exocomp/types"
import "embed"
import "io/fs"
import "path/filepath"

var Providers map[string]*types.Provider

//go:embed *.yaml
var filesystem embed.FS

func init() {

	Providers = make(map[string]*types.Provider)

	entries, err0 := fs.ReadDir(filesystem, ".")

	if err0 == nil {

		for _, entry := range entries {

			name := entry.Name()
			ext  := filepath.Ext(name)

			if ext == ".yaml" {

				data, err1 := filesystem.ReadFile(name)

				if err1 == nil {

					provider, err := types.ParseProvider(data)

					if err == nil {

						if provider.Model != "" {
							Providers[provider.Model] = provider
						}

					}

				}

			}

		}

	}

}
