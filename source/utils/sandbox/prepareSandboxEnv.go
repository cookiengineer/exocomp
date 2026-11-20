package sandbox

import "os"
import "path/filepath"

// prepareSandboxEnv prepares the private runtime environment for a sandboxed
// program. It returns the environment and the runtime paths that must be
// granted to the sandbox helper.
func prepareSandboxEnv(sandbox string) ([]string, []string, error) {

	base, err0 := sandboxRuntimeDir(sandbox)

	if err0 == nil {

		home := filepath.Join(base, "home")
		tmp := filepath.Join(base, "tmp")
		cache := filepath.Join(base, "cache")
		config := filepath.Join(base, "config")
		data := filepath.Join(base, "data")
		state := filepath.Join(base, "state")
		gopath := filepath.Join(base, "gopath")

		directories := []string{
			base,
			home,
			tmp,
			cache,
			config,
			data,
			state,
			gopath,
			filepath.Join(gopath, "pkg", "mod"),
		}

		err1 := prepareSandboxDirectories(directories)

		if err1 == nil {

			environment := os.Environ()

			environment = setEnv(environment, "HOME", home)
			environment = setEnv(environment, "XDG_CACHE_HOME", cache)
			environment = setEnv(environment, "XDG_CONFIG_HOME", config)
			environment = setEnv(environment, "XDG_DATA_HOME", data)
			environment = setEnv(environment, "XDG_STATE_HOME", state)
			environment = setEnv(environment, "XDG_RUNTIME_DIR", state)
			environment = setEnv(environment, "TMPDIR", tmp)
			environment = setEnv(environment, "TMP", tmp)
			environment = setEnv(environment, "TEMP", tmp)
			environment = setEnv(environment, "GOCACHE", cache)
			environment = setEnv(environment, "GOTMPDIR", tmp)
			environment = setEnv(environment, "GOPATH", gopath)
			environment = setEnv(environment, "GOMODCACHE", filepath.Join(gopath, "pkg", "mod"))
			environment = setEnv(environment, "GOENV", "off")
			environment = setEnv(environment, "GOTOOLCHAIN", "local")
			environment = unsetEnv(environment, "EXOCOMP_SANDBOX")
			environment = unsetEnv(environment, "SSH_AUTH_SOCK")
			environment = unsetEnv(environment, "SSH_AGENT_PID")
			environment = unsetEnv(environment, "GPG_AGENT_INFO")
			environment = unsetEnv(environment, "DBUS_SESSION_BUS_ADDRESS")
			environment = unsetEnv(environment, "KRB5CCNAME")
			environment = unsetEnv(environment, "DISPLAY")
			environment = unsetEnv(environment, "WAYLAND_DISPLAY")

			for _, name := range secretEnvironmentVariables {
				environment = unsetEnv(environment, name)
			}

			return environment, []string{base}, nil

		} else {
			return nil, nil, err1
		}

	} else {
		return nil, nil, err0
	}

}
