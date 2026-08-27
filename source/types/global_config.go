package types

import "os"
import "path/filepath"
import "sync"

var global_config_once sync.Once
var global_config *Config

func init() {

	global_config_once.Do(func() {

		home, err0 := os.UserHomeDir()

		if err0 == nil {

			path := filepath.Join(home, ".config", "exocomp", "config.yaml")

			data, err1 := os.ReadFile(path)

			if err1 == nil {

				config, err2 := ParseConfig(data)

				if err2 == nil {
					global_config = config
				}

			}

		}

	})

}

