package shell

import utils_sandbox "exocomp/utils/sandbox"
import utils_shell "exocomp/utils/shell"
import "context"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "slices"
import "strings"
import "time"
import pty "github.com/creack/pty"

func (tool *Shell) Open(program string, arguments []string, workdir string, columns int, rows int) (string, error) {

	program = strings.TrimSpace(program)

	if program == "" {
		return "", fmt.Errorf("shell.Open: Invalid parameter \"program\" is empty.")
	}

	if strings.ContainsAny(program, "/\\") == true {
		return "", fmt.Errorf("shell.Open: Invalid program \"%s\": Program name must not contain a path separator.", program)
	}

	if slices.Contains(tool.AllowedPrograms, program) == false {
		return "", fmt.Errorf("shell.Open: Invalid program \"%s\": Attempt to execute unallowed program.", program)
	}

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

}
