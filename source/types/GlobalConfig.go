package types

import "os"
import "path/filepath"

var GlobalConfig *Config

func init() {

	home, err0 := os.UserHomeDir()

	if err0 == nil {

		path := filepath.Join(home, ".config", "exocomp", "config.yaml")

		data, err1 := os.ReadFile(path)

		if err1 == nil {

			config, err2 := ParseConfig(data)

			if err2 == nil {
				GlobalConfig = config
			}

		}

	}

}

