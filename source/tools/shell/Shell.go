package shell

import "exocomp/schemas"
import utils_capabilities "exocomp/utils/capabilities"
import utils_sandbox "exocomp/utils/sandbox"
import utils_shell "exocomp/utils/shell"
import "context"
import "crypto/rand"
import "encoding/hex"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "slices"
import "sort"
import "strings"
import "sync"
import "time"

import pty "github.com/creack/pty"

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

func (tool *Shell) Call(method string, arguments map[string]any) (string, error) {

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

func (tool *Shell) activeSessions() int {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	count := 0

	for _, session := range tool.sessions {

		session.mutex.Lock()

		if session.Status == "starting" || session.Status == "running" {
			count++
		}

		session.mutex.Unlock()

	}

	return count

}

func (tool *Shell) callOpen(program string, arguments []string, workdir string, raw map[string]any) (string, error) {

	columns, err0 := utils_shell.ParseInt(raw, "columns", shellDefaultColumns)

	if err0 == nil {

		rows, err1 := utils_shell.ParseInt(raw, "rows", shellDefaultRows)

		if err1 == nil {
			return tool.Open(program, arguments, workdir, columns, rows)
		} else {
			return "", fmt.Errorf("shell.Open: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("shell.Open: %s", err0.Error())
	}

}

func (tool *Shell) cleanup() {

	now := time.Now()
	expired := make([]*shellSession, 0)

	tool.mutex.Lock()

	for session_id, session := range tool.sessions {

		session.mutex.Lock()
		status := session.Status
		last_activity := session.LastActivity
		started_at := session.StartedAt
		finished_at := session.FinishedAt
		session.mutex.Unlock()

		if (status == "starting" || status == "running") && (now.Sub(last_activity) > tool.IdleTimeout || now.Sub(started_at) > tool.Timeout) {

			delete(tool.sessions, session_id)
			expired = append(expired, session)

		} else if (status == "exited" || status == "failed" || status == "timed_out") && now.Sub(finished_at) > shellExitedRetention {

			delete(tool.sessions, session_id)
			expired = append(expired, session)

		}

	}

	tool.mutex.Unlock()

	for _, session := range expired {

		session.timeout()
		session.close(true)

	}

}

func (tool *Shell) Close(session_id string, is_force bool) (string, error) {

	tool.mutex.Lock()
	session, ok := tool.sessions[session_id]

	if ok == true {
		delete(tool.sessions, session_id)
	}

	tool.mutex.Unlock()

	if ok == true && session != nil {

		data, status, exit_code := session.close(is_force)
		total, _, _ := session.outputStats()

		header := fmt.Sprintf("shell.Close: Session \"%s\" closed with status \"%s\"%s, log %d bytes.", session_id, status, utils_shell.FormatExitCode(status, exit_code), total)

		if len(data) > 0 {
			return header + "\n" + string(data), nil
		} else {
			return header + "\n<no new output>", nil
		}

	} else {
		return "", fmt.Errorf("shell.Close: Session \"%s\" not found.", session_id)
	}

}

func (tool *Shell) getSession(session_id string) *shellSession {

	tool.mutex.Lock()
	session, ok := tool.sessions[session_id]
	tool.mutex.Unlock()

	if ok == true {
		return session
	} else {
		return nil
	}

}

func (tool *Shell) janitor() {

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {

		select {

		case <-tool.stop:
			return
		case <-ticker.C:
			tool.cleanup()

		}

	}

}

func (tool *Shell) List() (string, error) {

	tool.mutex.Lock()
	sessions := make([]*shellSession, 0, len(tool.sessions))

	for _, session := range tool.sessions {
		sessions = append(sessions, session)
	}

	tool.mutex.Unlock()

	sort.Slice(sessions, func(a int, b int) bool {
		return sessions[a].StartedAt.Before(sessions[b].StartedAt)
	})

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("shell.List: %d sessions.", len(sessions)))

	for _, session := range sessions {

		session.mutex.Lock()
		status := session.Status
		exit_code := session.ExitCode
		program := session.Program
		started_at := session.StartedAt
		last_activity := session.LastActivity
		log_size := session.logSize
		dropped := session.dropped
		session.mutex.Unlock()

		lines = append(lines, fmt.Sprintf("- Session: %s, Program: \"%s\", Status: %s%s, Started: %s, LastActivity: %s, Log: %d bytes, Dropped: %d bytes", session.ID, program, status, utils_shell.FormatExitCode(status, exit_code), started_at.Format(time.RFC3339), last_activity.Format(time.RFC3339), log_size, dropped))

	}

	return strings.Join(lines, "\n"), nil

}

func (tool *Shell) Log(session_id string, offset int, limit int) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if offset >= 0 {

			if limit >= 1 && limit <= shellMaxRead {

				data, total, truncated, err0 := session.readLog(int64(offset), limit)

				if err0 == nil {

					header := fmt.Sprintf("shell.Log: Session \"%s\" log %d bytes, offset %d, read %d bytes.", session_id, total, offset, len(data))

					if truncated == true {
						header += " Log is truncated."
					}

					if len(data) > 0 {
						return header + "\n" + string(data), nil
					} else {
						return header + "\n<no data>", nil
					}

				} else {
					return "", fmt.Errorf("shell.Log: %s", err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.Log: Invalid parameter \"limit\" must be between 1 and %d.", shellMaxRead)
			}

		} else {
			return "", fmt.Errorf("shell.Log: Invalid parameter \"offset\" must be zero or positive.")
		}

	} else {
		return "", fmt.Errorf("shell.Log: Session \"%s\" not found.", session_id)
	}

}

func (tool *Shell) newSessionID() string {

	for attempt := 0; attempt < 32; attempt++ {

		buffer := make([]byte, 4)
		_, err := rand.Read(buffer)

		if err == nil {

			session_id := "sh_" + hex.EncodeToString(buffer)

			tool.mutex.Lock()
			_, exists := tool.sessions[session_id]
			tool.mutex.Unlock()

			if exists == false {
				return session_id
			}

		} else {
			return ""
		}

	}

	return ""

}

func (tool *Shell) Open(program string, arguments []string, workdir string, columns int, rows int) (string, error) {

	program = strings.TrimSpace(program)

	if program != "" {

		if strings.ContainsAny(program, "/\\") == false {

			if slices.Contains(tool.AllowedPrograms, program) == true {

				if workdir == "" {
					workdir = "."
				}

				if columns == 0 {
					columns = shellDefaultColumns
				}

				if rows == 0 {
					rows = shellDefaultRows
				}

				if columns < shellMinColumns || columns > shellMaxColumns {
					return "", fmt.Errorf("shell.Open: Invalid parameter \"columns\" must be between %d and %d.", shellMinColumns, shellMaxColumns)
				}

				if rows < shellMinRows || rows > shellMaxRows {
					return "", fmt.Errorf("shell.Open: Invalid parameter \"rows\" must be between %d and %d.", shellMinRows, shellMaxRows)
				}

				if utils_shell.ArgumentBytes(arguments) > shellMaxArguments {
					return "", fmt.Errorf("shell.Open: Invalid parameter \"arguments\" exceeds the %d byte limit.", shellMaxArguments)
				}

				resolved_arguments, err_args := sanitizeShellArguments(tool.Sandbox, tool.Capabilities, arguments)

				if err_args != nil {
					return "", fmt.Errorf("shell.Open: %s", err_args.Error())
				}

				resolved_workdir, err0 := utils_shell.ResolveWorkdir(tool.Sandbox, workdir)

				if err0 == nil {

					if tool.activeSessions() >= tool.MaxSessions {
						return "", fmt.Errorf("shell.Open: Session limit reached (%d active sessions). Close a session first.", tool.MaxSessions)
					}

					_, err1 := exec.LookPath(program)

					if err1 == nil {

						log_dir, err2 := utils_sandbox.ResolvePath(tool.Sandbox, filepath.Join(".exocomp", "shell"))

						if err2 == nil {

							mkdir_err := os.MkdirAll(log_dir, 0755)

							if mkdir_err == nil {

								session_id := tool.newSessionID()

								if session_id != "" {

									log_path := filepath.Join(log_dir, session_id+".log")
									log_file, err3 := os.OpenFile(log_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)

									if err3 == nil {

										ctx, cancel := context.WithCancel(context.Background())
										read_paths, write_paths := tool.Capabilities.SandboxPaths()

										cmd, err4 := utils_sandbox.CommandContextWithPermissions(ctx, tool.Sandbox, program, resolved_arguments, utils_sandbox.Permissions{
											ReadPaths:  read_paths,
											WritePaths: write_paths,
											Network:    tool.Capabilities.AllowsNetworkEgress(),
										})

										if err4 == nil {

											cmd.Dir = resolved_workdir
											cmd.Env = utils_shell.SetEnv(cmd.Env, "TERM", "xterm-256color")
											cmd.Env = utils_shell.SetEnv(cmd.Env, "PAGER", "cat")
											cmd.Env = utils_shell.SetEnv(cmd.Env, "GIT_PAGER", "cat")
											cmd.Env = utils_shell.SetEnv(cmd.Env, "SYSTEMD_PAGER", "cat")

											utils_shell.PrepareProcess(cmd)

											master, err5 := pty.StartWithSize(cmd, &pty.Winsize{
												Cols: uint16(columns),
												Rows: uint16(rows),
											})

											if err5 == nil {

												now := time.Now()
												session := &shellSession{
													ID:           session_id,
													Program:      program,
													Arguments:    append([]string{}, resolved_arguments...),
													Workdir:      workdir,
													Columns:      columns,
													Rows:         rows,
													Status:       "running",
													ExitCode:     -1,
													StartedAt:    now,
													LastActivity: now,
													cmd:          cmd,
													pty:          master,
													cancel:       cancel,
													logFile:      log_file,
													logPath:      log_path,
													notify:       make(chan struct{}, 1),
													done:         make(chan struct{}),
												}

												tool.mutex.Lock()
												tool.sessions[session_id] = session
												tool.mutex.Unlock()

												go session.run()

												initial, status, exit_code := session.read(shellOpenWait, shellDefaultRead)

												lines := make([]string, 0)
												lines = append(lines, fmt.Sprintf("shell.Open: Session \"%s\" started program \"%s\" arguments %s in workdir \"%s\".", session_id, program, utils_shell.FormatArguments(resolved_arguments), workdir))
												lines = append(lines, fmt.Sprintf("shell.Open: Status: %s%s.", status, utils_shell.FormatExitCode(status, exit_code)))
												lines = append(lines, "shell.Open: Output:")

												if len(initial) > 0 {
													lines = append(lines, string(initial))
												} else {
													lines = append(lines, "<none>")
												}

												return strings.Join(lines, "\n"), nil

											} else {

												cancel()
												log_file.Close()
												os.Remove(log_path)
												return "", fmt.Errorf("shell.Open: Cannot start program \"%s\".", program)

											}

										} else {

											cancel()
											log_file.Close()
											os.Remove(log_path)
											return "", fmt.Errorf("shell.Open: %s", err4.Error())

										}

									} else {
										return "", fmt.Errorf("shell.Open: Cannot create session log.")
									}

								} else {
									return "", fmt.Errorf("shell.Open: Cannot create session identifier.")
								}

							} else {
								return "", fmt.Errorf("shell.Open: Cannot create session directory.")
							}

						} else {
							return "", fmt.Errorf("shell.Open: Cannot create session directory.")
						}

					} else {
						return "", fmt.Errorf("shell.Open: Invalid program \"%s\": Program doesn't exist.", program)
					}

				} else {
					return "", fmt.Errorf("shell.Open: %s", err0.Error())
				}

			} else {
				return "", fmt.Errorf("shell.Open: Invalid program \"%s\": Attempt to execute unallowed program.", program)
			}
		} else {
			return "", fmt.Errorf("shell.Open: Invalid program \"%s\": Program name must not contain a path separator.", program)
		}
	} else {
		return "", fmt.Errorf("shell.Open: Invalid parameter \"program\" is empty.")
	}

}

func (tool *Shell) Read(session_id string, wait_ms int, limit int) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if wait_ms >= 0 && wait_ms <= 30000 {

			if limit >= 1 && limit <= shellMaxRead {

				data, status, exit_code := session.read(time.Duration(wait_ms)*time.Millisecond, limit)
				dropped, total, truncated := session.outputStats()
				pending := session.pendingOutput()

				header := fmt.Sprintf("shell.Read: Session \"%s\" status \"%s\" read %d bytes (total %d, dropped %d)%s.", session_id, status, len(data), total, dropped, utils_shell.FormatExitCode(status, exit_code))

				if truncated == true || pending > 0 || dropped > 0 {
					header += " Output was truncated; use shell.Log."
				}

				if len(data) > 0 {
					return header + "\n" + string(data), nil
				} else {
					return header + "\n<no new output>", nil
				}

			} else {
				return "", fmt.Errorf("shell.Read: Invalid parameter \"limit\" must be between 1 and %d.", shellMaxRead)
			}

		} else {
			return "", fmt.Errorf("shell.Read: Invalid parameter \"wait_ms\" must be between 0 and 30000.")
		}

	} else {
		return "", fmt.Errorf("shell.Read: Session \"%s\" not found.", session_id)
	}

}

func (tool *Shell) Resize(session_id string, columns int, rows int) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if columns >= shellMinColumns && columns <= shellMaxColumns {

			if rows >= shellMinRows && rows <= shellMaxRows {

				err0 := session.resize(columns, rows)

				if err0 == nil {
					return fmt.Sprintf("shell.Resize: Session \"%s\" resized to %dx%d.", session_id, columns, rows), nil
				} else {
					return "", fmt.Errorf("shell.Resize: Cannot resize session \"%s\".", session_id)
				}

			} else {
				return "", fmt.Errorf("shell.Resize: Invalid parameter \"rows\" must be between %d and %d.", shellMinRows, shellMaxRows)
			}

		} else {
			return "", fmt.Errorf("shell.Resize: Invalid parameter \"columns\" must be between %d and %d.", shellMinColumns, shellMaxColumns)
		}

	} else {
		return "", fmt.Errorf("shell.Resize: Session \"%s\" not found.", session_id)
	}

}

func (tool *Shell) Signal(session_id string, name string) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		signal, err0 := session.signal(name)

		if err0 == nil {

			status, exit_code := session.status()

			return fmt.Sprintf("shell.Signal: Session \"%s\" sent SIG%s, status \"%s\"%s.", session_id, signal, status, utils_shell.FormatExitCode(status, exit_code)), nil

		} else {
			return "", fmt.Errorf("shell.Signal: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("shell.Signal: Session \"%s\" not found.", session_id)
	}

}

func (tool *Shell) Stat(session_id string) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		session.mutex.Lock()
		program := session.Program
		arguments := append([]string{}, session.Arguments...)
		workdir := session.Workdir
		status := session.Status
		exit_code := session.ExitCode
		started_at := session.StartedAt
		finished_at := session.FinishedAt
		last_activity := session.LastActivity
		columns := session.Columns
		rows := session.Rows
		log_size := session.logSize
		dropped := session.dropped
		log_truncated := session.logTruncated
		session.mutex.Unlock()

		finished := "-"

		if finished_at.IsZero() == false {
			finished = finished_at.Format(time.RFC3339)
		}

		truncated := ""

		if log_truncated == true {
			truncated = ", LogTruncated: true"
		}

		return strings.Join([]string{
			fmt.Sprintf("shell.Stat: Session \"%s\": program \"%s\", arguments %s, status \"%s\"%s.", session_id, program, utils_shell.FormatArguments(arguments), status, utils_shell.FormatExitCode(status, exit_code)),
			fmt.Sprintf("shell.Stat: Workdir: \"%s\", TTY: %dx%d, Started: %s, Finished: %s, LastActivity: %s, Log: %d bytes, Dropped: %d bytes%s.", workdir, columns, rows, started_at.Format(time.RFC3339), finished, last_activity.Format(time.RFC3339), log_size, dropped, truncated),
		}, "\n"), nil

	} else {
		return "", fmt.Errorf("shell.Stat: Session \"%s\" not found.", session_id)
	}

}

func (tool *Shell) Write(session_id string, data string) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if len(data) <= shellMaxWrite {

			written, err0 := session.write([]byte(data))

			if err0 == nil {

				status, exit_code := session.status()

				return fmt.Sprintf("shell.Write: Session \"%s\" wrote %d bytes, status \"%s\"%s.", session_id, written, status, utils_shell.FormatExitCode(status, exit_code)), nil

			} else {
				return "", fmt.Errorf("shell.Write: Cannot write to session \"%s\".", session_id)
			}

		} else {
			return "", fmt.Errorf("shell.Write: Invalid parameter \"data\" exceeds the %d byte limit.", shellMaxWrite)
		}

	} else {
		return "", fmt.Errorf("shell.Write: Session \"%s\" not found.", session_id)
	}

}
