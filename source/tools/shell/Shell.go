package shell

import "exocomp/schemas"
import utils_capabilities "exocomp/utils/capabilities"
import utils_shell "exocomp/utils/shell"
import "fmt"
import "slices"
import "sync"
import "time"

const (
	shellDefaultColumns  = 80
	shellDefaultRows     = 24
	shellMinColumns      = 20
	shellMaxColumns      = 500
	shellMinRows         = 5
	shellMaxRows         = 200
	shellDefaultRead     = 8192
	shellMaxRead         = 65536
	shellMaxWrite        = 65536
	shellMaxArguments    = 65536
	shellMaxBuffer       = 1024 * 1024
	shellMaxLog          = 16 * 1024 * 1024
	shellOpenWait        = 250 * time.Millisecond
	shellCloseGrace      = 2 * time.Second
	shellCloseWait       = 5 * time.Second
	shellExitedRetention = 5 * time.Minute
)

type Shell struct {
	Methods         []string
	Playground      string
	Sandbox         string
	AllowedPrograms []string
	Capabilities    *utils_capabilities.Set
	Timeout         time.Duration
	IdleTimeout     time.Duration
	MaxSessions     int

	mutex    *sync.Mutex
	sessions map[string]*shellSession
	stop     chan struct{}
	stopOnce *sync.Once
}

func NewShell(methods []string, playground string, sandbox string, allowed_programs []string) *Shell {

	return NewShellWithCapabilities(methods, playground, sandbox, allowed_programs, nil)

}

func NewShellWithCapabilities(
	methods []string,
	playground string,
	sandbox string,
	allowed_programs []string,
	capabilities *utils_capabilities.Set,
) *Shell {

	tool := &Shell{
		Methods:         methods,
		Playground:      playground,
		Sandbox:         sandbox,
		AllowedPrograms: allowed_programs,
		Capabilities:    capabilities,
		Timeout:         30 * time.Minute,
		IdleTimeout:     5 * time.Minute,
		MaxSessions:     4,
		mutex:           &sync.Mutex{},
		sessions:        make(map[string]*shellSession),
		stop:            make(chan struct{}),
		stopOnce:        &sync.Once{},
	}

	go tool.janitor()

	return tool

}

func (tool *Shell) Name() string {
	return "shell"
}

func (tool *Shell) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Open" {

			program, ok1 := arguments["program"].(string)

			if ok1 == true {

				args, err0 := utils_shell.ParseArguments(arguments["arguments"])

				if err0 == nil {

					workdir := "."
					raw_workdir, ok2 := arguments["workdir"]

					if ok2 == true {

						tmp, ok3 := raw_workdir.(string)

						if ok3 == true {
							workdir = tmp
							return tool.callOpen(program, args, workdir, arguments)
						} else {
							return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"workdir\" is not a string.")
						}

					} else {
						return tool.callOpen(program, args, workdir, arguments)
					}

				} else {
					return "", fmt.Errorf("shell.%s: %s", method, err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"program\" is not a string.")
			}

		} else if method == "Write" {

			session_id, ok1 := arguments["session_id"].(string)
			data, ok2 := arguments["data"].(string)

			if ok1 == true && ok2 == true {
				return tool.Write(session_id, data)
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"data\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Read" {

			session_id, ok := arguments["session_id"].(string)

			if ok == true {

				wait_ms, err0 := utils_shell.ParseInt(arguments, "wait_ms", 0)

				if err0 == nil {

					limit, err1 := utils_shell.ParseInt(arguments, "limit", shellDefaultRead)

					if err1 == nil {
						return tool.Read(session_id, wait_ms, limit)
					} else {
						return "", fmt.Errorf("shell.%s: %s", method, err1.Error())
					}

				} else {
					return "", fmt.Errorf("shell.%s: %s", method, err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			}

		} else if method == "Resize" {

			session_id, ok := arguments["session_id"].(string)

			if ok == true {

				columns, err0 := utils_shell.ParseInt(arguments, "columns", 0)

				if err0 == nil {

					rows, err1 := utils_shell.ParseInt(arguments, "rows", 0)

					if err1 == nil {
						return tool.Resize(session_id, columns, rows)
					} else {
						return "", fmt.Errorf("shell.%s: %s", method, err1.Error())
					}

				} else {
					return "", fmt.Errorf("shell.%s: %s", method, err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			}

		} else if method == "Signal" {

			session_id, ok1 := arguments["session_id"].(string)
			signal, ok2 := arguments["signal"].(string)

			if ok1 == true && ok2 == true {
				return tool.Signal(session_id, signal)
			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"signal\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameters.")
			}

		} else if method == "Close" {

			session_id, ok := arguments["session_id"].(string)

			if ok == true {

				force := false
				raw_force, ok2 := arguments["force"]

				if ok2 == true {

					tmp, ok3 := raw_force.(bool)

					if ok3 == true {
						force = tmp
						return tool.Close(session_id, force)
					} else {
						return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"force\" is not a boolean.")
					}

				} else {
					return tool.Close(session_id, force)
				}

			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			}

		} else if method == "List" {

			return tool.List()

		} else if method == "Stat" {

			session_id, ok := arguments["session_id"].(string)

			if ok == true {
				return tool.Stat(session_id)
			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			}

		} else if method == "Log" {

			session_id, ok := arguments["session_id"].(string)

			if ok == true {

				offset, err0 := utils_shell.ParseInt(arguments, "offset", 0)

				if err0 == nil {

					limit, err1 := utils_shell.ParseInt(arguments, "limit", shellDefaultRead)

					if err1 == nil {
						return tool.Log(session_id, offset, limit)
					} else {
						return "", fmt.Errorf("shell.%s: %s", method, err1.Error())
					}

				} else {
					return "", fmt.Errorf("shell.%s: %s", method, err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.%s: %s", method, "Invalid parameter \"session_id\" is not a string.")
			}

		} else {
			return "", fmt.Errorf("shell.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("shell.%s: Method not allowed.", method)
	}

}

func (tool *Shell) Destroy() {

	if tool != nil {

		if tool.stopOnce != nil {
			tool.stopOnce.Do(func() {
				close(tool.stop)
			})
		}

		tool.mutex.Lock()
		sessions := make([]*shellSession, 0, len(tool.sessions))

		for session_id, session := range tool.sessions {
			delete(tool.sessions, session_id)
			sessions = append(sessions, session)
		}

		tool.mutex.Unlock()

		for _, session := range sessions {
			session.close(true)
		}

	}

}

func (tool *Shell) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Shell) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *Shell) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Shell) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range ShellSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
