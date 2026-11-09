package sandbox

import "os"
import "path/filepath"

func prepareSandboxEnv(sandbox string) ([]string, error) {

	base := filepath.Join(sandbox, ".exocomp", sandboxDirectory())
	home := filepath.Join(base, "home")
	tmp := filepath.Join(base, "tmp")
	cache := filepath.Join(base, "cache")
	config := filepath.Join(base, "config")
	data := filepath.Join(base, "data")
	state := filepath.Join(base, "state")
	gopath := filepath.Join(base, "gopath")

	directories := []string{
		home,
		tmp,
		cache,
		config,
		data,
		state,
		gopath,
		filepath.Join(gopath, "pkg", "mod"),
	}

	err0 := prepareSandboxDirectories(directories)

	if err0 == nil {

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

		return environment, nil

	} else {
		return nil, err0
	}

}
