package shell

import "exocomp/tools/testutil"
import "os"
import "path/filepath"
import "strings"
import "time"
import "testing"

func newTestShell(t *testing.T) (*Shell, string) {

	testutil.RequireSandbox(t)

	fake_shell := testutil.BuildFakeShell(t)
	path := filepath.Dir(fake_shell) + string(os.PathListSeparator) + os.Getenv("PATH")
	t.Setenv("PATH", path)

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "work")

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Fatalf("Cannot create sandbox: %v", err0)
	}

	tool := NewShell(
		[]string{"Open", "Write", "Read", "Resize", "Signal", "Close", "List", "Stat", "Log"},
		playground,
		sandbox,
		[]string{filepath.Base(fake_shell)},
	)

	t.Cleanup(tool.Destroy)

	return tool, sandbox

}

func TestShellSchema(t *testing.T) {

	expected := []string{
		"shell.Open",
		"shell.Write",
		"shell.Read",
		"shell.Resize",
		"shell.Signal",
		"shell.Close",
		"shell.List",
		"shell.Stat",
		"shell.Log",
	}

	if len(ShellSchema) != len(expected) {
		t.Fatalf("Expected %d shell schemas, got %d", len(expected), len(ShellSchema))
	}

	for index, name := range expected {

		if ShellSchema[index].Function.Name != name {
			t.Errorf("Expected schema %q, got %q", name, ShellSchema[index].Function.Name)
		}

	}

}

func TestShell_OpenReadWriteClose(t *testing.T) {

	tool, _ := newTestShell(t)

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		if strings.Contains(result0, "prompt>") == false {
			t.Errorf("Expected initial prompt, got %s", result0)
		}

		session_id := extractShellSessionID(t, result0)

		result1, err1 := tool.Write(session_id, "echo hello\n")

		if err1 == nil {

			if strings.Contains(result1, "wrote 11 bytes") == false {
				t.Errorf("Expected write report, got %s", result1)
			}

			result2, err2 := tool.Read(session_id, 1000, shellDefaultRead)

			if err2 == nil {

				if strings.Contains(result2, "hello") == false {
					t.Errorf("Expected echo output, got %s", result2)
				}

				result3, err3 := tool.Write(session_id, "exit\n")

				if err3 == nil {

					if strings.Contains(result3, "wrote 5 bytes") == false {
						t.Errorf("Expected exit write report, got %s", result3)
					}

					result4 := readShellUntil(t, tool, session_id, "bye")

					if strings.Contains(result4, "bye") == false {
						t.Errorf("Expected exit output, got %s", result4)
					}

					result5, err5 := tool.Close(session_id, false)

					if err5 == nil {

						if strings.Contains(result5, "status \"exited\"") == false {
							t.Errorf("Expected exited status, got %s", result5)
						}

					} else {
						t.Errorf("Expected %v to be nil", err5)
					}

				} else {
					t.Errorf("Expected %v to be nil", err3)
				}

			} else {
				t.Errorf("Expected %v to be nil", err2)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_Resize(t *testing.T) {

	tool, _ := newTestShell(t)

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		session_id := extractShellSessionID(t, result0)

		result1, err1 := tool.Resize(session_id, 100, 40)

		if err1 == nil {

			if strings.Contains(result1, "100x40") == false {
				t.Errorf("Expected resize report, got %s", result1)
			}

			_, err2 := tool.Write(session_id, "size\n")

			if err2 == nil {

				result3 := readShellUntil(t, tool, session_id, "size: 100x40")

				if strings.Contains(result3, "size: 100x40") == false {
					t.Errorf("Expected resized terminal, got %s", result3)
				}

			} else {
				t.Errorf("Expected %v to be nil", err2)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

		_, _ = tool.Close(session_id, true)

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_AnsiStripped(t *testing.T) {

	tool, _ := newTestShell(t)

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		session_id := extractShellSessionID(t, result0)

		_, err1 := tool.Write(session_id, "ansi\n")

		if err1 == nil {

			result2 := readShellUntil(t, tool, session_id, "red")

			if strings.Contains(result2, "red") == false {
				t.Errorf("Expected ANSI content, got %s", result2)
			}

			if strings.Contains(result2, "\x1b") == true {
				t.Errorf("Expected ANSI escapes to be stripped, got %s", result2)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

		_, _ = tool.Close(session_id, true)

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_LogAndTruncation(t *testing.T) {

	tool, _ := newTestShell(t)

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		session_id := extractShellSessionID(t, result0)

		_, err1 := tool.Write(session_id, "big\n")

		if err1 == nil {

			result2 := readShellUntilLimit(t, tool, session_id, "Output was truncated", 128)

			if strings.Contains(result2, "Output was truncated") == false {
				t.Errorf("Expected truncation marker, got %s", result2)
			}

			result3, err3 := tool.Log(session_id, 0, 4096)

			if err3 == nil {

				if strings.Count(result3, "A") < 1000 {
					t.Errorf("Expected transcript content, got %s", result3)
				}

				result4, err4 := tool.Log(session_id, 999999, 64)

				if err4 == nil {

					if strings.Contains(result4, "<no data>") == false {
						t.Errorf("Expected empty log segment, got %s", result4)
					}

				} else {
					t.Errorf("Expected %v to be nil", err4)
				}

			} else {
				t.Errorf("Expected %v to be nil", err3)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

		_, _ = tool.Close(session_id, true)

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_Signal(t *testing.T) {

	tool, _ := newTestShell(t)

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		session_id := extractShellSessionID(t, result0)

		_, err1 := tool.Write(session_id, "wait\n")

		if err1 == nil {

			_ = readShellUntil(t, tool, session_id, "waiting")

			result3, err3 := tool.Signal(session_id, "INT")

			if err3 == nil {

				if strings.Contains(result3, "SIGINT") == false {
					t.Errorf("Expected SIGINT report, got %s", result3)
				}

				result4 := readShellUntilNotRunning(t, tool, session_id)

				if strings.Contains(result4, "status \"running\"") == true {
					t.Errorf("Expected session to stop, got %s", result4)
				}

			} else {
				t.Errorf("Expected %v to be nil", err3)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

		_, _ = tool.Close(session_id, true)

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_ListStatAndCloseMissing(t *testing.T) {

	tool, _ := newTestShell(t)

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		session_id := extractShellSessionID(t, result0)

		result1, err1 := tool.List()

		if err1 == nil {

			if strings.Contains(result1, session_id) == false {
				t.Errorf("Expected session %s in list, got %s", session_id, result1)
			}

			result2, err2 := tool.Stat(session_id)

			if err2 == nil {

				if strings.Contains(result2, "prompt") == true {
					t.Errorf("Expected only session metadata, got %s", result2)
				}

				if strings.Contains(result2, "Status: \"running\"") == false && strings.Contains(result2, "status \"running\"") == false {
					t.Errorf("Expected running status, got %s", result2)
				}

			} else {
				t.Errorf("Expected %v to be nil", err2)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

		_, err3 := tool.Close(session_id, true)

		if err3 == nil {

			_, err4 := tool.Close(session_id, true)

			if err4 != nil {

				if strings.Contains(err4.Error(), "not found") == false {
					t.Errorf("Expected not-found error, got %v", err4)
				}

			} else {
				t.Errorf("Expected %v to be not nil", err4)
			}

		} else {
			t.Errorf("Expected %v to be nil", err3)
		}

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_Validation(t *testing.T) {

	tool, sandbox := newTestShell(t)

	_, err0 := tool.Open("definitely-not-allowed", []string{}, ".", 80, 24)

	if err0 != nil {

		if strings.Contains(err0.Error(), "unallowed program") == false {
			t.Errorf("Expected unallowed program error, got %v", err0)
		}

	} else {
		t.Errorf("Expected %v to be not nil", err0)
	}

	_, err1 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, "../outside", 80, 24)

	if err1 != nil {

		if strings.Contains(err1.Error(), "sandbox") == false && strings.Contains(err1.Error(), "workdir") == false {
			t.Errorf("Expected workdir error, got %v", err1)
		}

	} else {
		t.Errorf("Expected %v to be not nil", err1)
	}

	_, err2 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{"/etc/passwd"}, ".", 80, 24)

	if err2 != nil {

		if strings.Contains(err2.Error(), "escape sandbox") == false {
			t.Errorf("Expected argument sandbox error, got %v", err2)
		}

	} else {
		t.Errorf("Expected %v to be not nil", err2)
	}

	mkdir_err := os.MkdirAll(filepath.Join(sandbox, "sub"), 0755)

	if mkdir_err == nil {

		result3, open_err := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, "./sub", 80, 24)

		if open_err == nil {

			session_id := extractShellSessionID(t, result3)
			_, _ = tool.Close(session_id, true)

		} else {
			t.Errorf("Expected %v to be nil", open_err)
		}

	} else {
		t.Errorf("Expected %v to be nil", mkdir_err)
	}

}

func TestShell_SessionLimit(t *testing.T) {

	tool, _ := newTestShell(t)
	tool.MaxSessions = 1

	result0, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		session_id := extractShellSessionID(t, result0)

		_, err1 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

		if err1 != nil {

			if strings.Contains(err1.Error(), "Session limit reached") == false {
				t.Errorf("Expected session limit error, got %v", err1)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err1)
		}

		_, _ = tool.Close(session_id, true)

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_IdleTimeout(t *testing.T) {

	tool, _ := newTestShell(t)
	tool.IdleTimeout = 50 * time.Millisecond
	tool.Timeout = 2 * time.Second

	_, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		time.Sleep(100 * time.Millisecond)
		tool.cleanup()

		result1, err1 := tool.List()

		if err1 == nil {

			if strings.Contains(result1, "0 sessions") == false {
				t.Errorf("Expected expired session to be removed, got %s", result1)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func TestShell_Destroy(t *testing.T) {

	tool, _ := newTestShell(t)

	_, err0 := tool.Open(filepath.Base(testutil.BuildFakeShell(t)), []string{}, ".", 80, 24)

	if err0 == nil {

		tool.Destroy()

		result1, err1 := tool.List()

		if err1 == nil {

			if strings.Contains(result1, "0 sessions") == false {
				t.Errorf("Expected no sessions after destroy, got %s", result1)
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

}

func readShellUntil(t *testing.T, tool *Shell, session_id string, needle string) string {

	t.Helper()

	return readShellUntilLimit(t, tool, session_id, needle, shellDefaultRead)

}

func readShellUntilLimit(t *testing.T, tool *Shell, session_id string, needle string, limit int) string {

	t.Helper()

	output := ""
	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {

		result, err := tool.Read(session_id, 250, limit)

		if err == nil {

			output += result

			if strings.Contains(output, needle) == true {
				return output
			}

		} else {
			t.Fatalf("Expected %v to be nil", err)
		}

	}

	t.Fatalf("Expected %q in shell output, got %s", needle, output)

	return output

}

func readShellUntilNotRunning(t *testing.T, tool *Shell, session_id string) string {

	t.Helper()

	output := ""
	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {

		result, err := tool.Read(session_id, 250, shellDefaultRead)

		if err == nil {

			output += result

			if strings.Contains(result, "status \"running\"") == false && strings.Contains(result, "status \"starting\"") == false {
				return output
			}

		} else {
			t.Fatalf("Expected %v to be nil", err)
		}

	}

	t.Fatalf("Expected session to stop, got %s", output)

	return output

}

func extractShellSessionID(t *testing.T, result string) string {

	t.Helper()

	prefix := "shell.Open: Session \""

	if strings.HasPrefix(result, prefix) == true {

		rest := strings.TrimPrefix(result, prefix)
		index := strings.Index(rest, "\"")

		if index > 0 {
			return rest[0:index]
		} else {
			t.Fatalf("Cannot parse session identifier from %s", result)
			return ""
		}

	} else {
		t.Fatalf("Cannot parse session identifier from %s", result)
		return ""
	}

}
