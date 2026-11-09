package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestPrograms_Execute(t *testing.T) {

	requireLandlock(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"cat", "ls", "pwd"})

	err0 := os.MkdirAll(filepath.Join(sandbox, "probe"), 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if tool != nil {

		result1, err1 := tool.Execute("ls", []string{"-la", "./probe"})
		result2, err2 := tool.Execute("pwd", []string{})
		result3, err3 := tool.Execute("cat", []string{"./../../../file.txt"})
		result4, err4 := tool.Execute("cat", []string{"/etc/passwd"})
		result5, err5 := tool.Execute("nc", []string{"192.168.123.123", "1337"})

		lines1 := strings.Split(result1, "\n")

		if len(lines1) == 5 {

			if lines1[0] != "programs.Execute: ls -la ./probe" {
				t.Errorf("Expected \"%s\" to be executed", "ls -la ./probe")
			}

			if lines1[1] != "total 0" {
				t.Errorf("Expected %d files", 0)
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines1), 5)
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		lines2 := strings.Split(result2, "\n")

		if len(lines2) == 3 {

			if lines2[0] != "programs.Execute: pwd" {
				t.Errorf("Expected \"%s\" to be executed", "pwd")
			}

			if lines2[1] != sandbox {
				t.Errorf("Expected \"%s\" to be \"%s\"", lines2[1], sandbox)
			}

			if lines2[2] != "" {
				t.Errorf("Expected \"%s\" to be empty", lines2[2])
			}

		} else {
			t.Errorf("Expected %d lines to be %d", len(lines2), 3)
		}

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if result3 != "" {
			t.Errorf("Expected %s to be empty", result3)
		}

		if err3 != nil {

			if strings.Contains(err3.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err3)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err3)
		}

		if result4 != "" {
			t.Errorf("Expected %s to be empty", result4)
		}

		if err4 != nil {

			if strings.Contains(err4.Error(), "Attempt to escape sandbox") == false {
				t.Errorf("Expected %v to detect attempt to escape sandbox", err4)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err4)
		}

		if result5 != "" {
			t.Errorf("Expected %s to be empty", result5)
		}

		if err5 != nil {

			if strings.Contains(err5.Error(), "Attempt to execute unallowed program") == false {
				t.Errorf("Expected %v to detect attempt to execute unallowed program", err5)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err5)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_ExecuteWithoutPermission(t *testing.T) {

	requireLandlock(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"ls", "pwd"})

	// Folder with no execution rights
	err0 := os.MkdirAll(sandbox, 0644)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if tool != nil {

		result1, err1 := tool.Execute("ls", []string{"-la"})
		result2, err2 := tool.Execute("pwd", []string{})

		if result1 != "" {
			t.Errorf("Expected %s to be empty", result1)
		}

		if err1 != nil {

			if strings.Contains(err1.Error(), "Permission denied") == false {
				t.Errorf("Expected %v to detect denied permission", err1)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err1)
		}

		if result2 != "" {
			t.Errorf("Expected %s to be empty", result2)
		}

		if err2 != nil {

			if strings.Contains(err2.Error(), "Permission denied") == false {
				t.Errorf("Expected %v to detect denied permission", err2)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err2)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_ExecuteWithoutProgram(t *testing.T) {
	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"doesntexist"})

	// Folder with execution rights
	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if tool != nil {

		result1, err1 := tool.Execute("doesntexist", []string{})

		if result1 != "" {
			t.Errorf("Expected %s to be empty", result1)
		}

		if err1 != nil {

			if strings.Contains(err1.Error(), "Program doesn't exist") == false {
				t.Errorf("Expected %v to detect not existing program", err1)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err1)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"cat", "ls", "pwd"})

	if tool != nil {

		result1, err1 := tool.List()

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		lines := strings.Split(result1, "\n")

		if len(lines) != 4 {
			t.Errorf("Expected %d lines, got %d: %s", 4, len(lines), result1)
		}

		if lines[0] != "programs.List:" {
			t.Errorf("Expected header %q, got %q", "programs.List:", lines[0])
		}

		if strings.Contains(lines[1], "Name: cat") == false {
			t.Errorf("Expected %q to contain %q", lines[1], "Name: cat")
		}

		if strings.Contains(lines[2], "Name: ls") == false {
			t.Errorf("Expected %q to contain %q", lines[2], "Name: ls")
		}

		if strings.Contains(lines[3], "Name: pwd") == false {
			t.Errorf("Expected %q to contain %q", lines[3], "Name: pwd")
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_Stat(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"cat", "ls", "pwd"})

	if tool != nil {

		result1, err1 := tool.Stat("cat")
		result2, err2 := tool.Stat("nc")

		if result1 == "" {
			t.Errorf("Expected resolved path for %q", "cat")
		}

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if result2 != "" {
			t.Errorf("Expected %s to be empty", result2)
		}

		if err2 != nil {

			if strings.Contains(err2.Error(), "Attempt to lookup unallowed program") == false {
				t.Errorf("Expected %v to detect unallowed program", err2)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err2)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_StatWithoutProgram(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"doesntexist"})

	if tool != nil {

		result1, err1 := tool.Stat("doesntexist")

		if result1 != "" {
			t.Errorf("Expected %s to be empty", result1)
		}

		if err1 != nil {

			if strings.Contains(err1.Error(), "Program doesn't exist") == false {
				t.Errorf("Expected %v to detect not existing program", err1)
			}

		} else {
			t.Errorf("Expected %v to be not nil", err1)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_Call_OptionalArguments(t *testing.T) {

	requireLandlock(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox := filepath.Join(playground, "programs")
	tool := NewPrograms([]string{"List", "Execute", "Stat"}, playground, sandbox, []string{"pwd"})

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if tool != nil {

		result1, err1 := tool.Call("Execute", map[string]interface{}{
			"program": "pwd",
		})

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if !strings.HasPrefix(result1, "programs.Execute: pwd") {
			t.Errorf("Expected pwd to be executed without arguments, got %s", result1)
		}

		result2, err2 := tool.Call("Execute", map[string]interface{}{
			"program":   "pwd",
			"arguments": []interface{}{},
		})

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if !strings.HasPrefix(result2, "programs.Execute: pwd") {
			t.Errorf("Expected pwd to be executed with empty arguments, got %s", result2)
		}

		_, err3 := tool.Call("Execute", map[string]interface{}{
			"program":   "pwd",
			"arguments": "not-an-array",
		})

		if err3 == nil {
			t.Errorf("Expected a non-nil error")
		} else if strings.Contains(err3.Error(), "not an array of strings") == false {
			t.Errorf("Expected array error, got %v", err3)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_ExecuteSymlinks(t *testing.T) {

	requireLandlock(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-symlink-*")
	sandbox := filepath.Join(playground, "programs")
	outside := filepath.Join(playground, "outside")

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	err1 := os.MkdirAll(outside, 0755)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	err2 := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("TOP SECRET"), 0644)

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	err3 := os.WriteFile(filepath.Join(sandbox, "normal.txt"), []byte("normal content"), 0644)

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	createTestSymlink(t, filepath.Join(outside, "secret.txt"), filepath.Join(sandbox, "leak"))
	createTestSymlink(t, outside, filepath.Join(sandbox, "leakdir"))

	tool := NewPrograms([]string{"Execute", "List", "Stat"}, playground, sandbox, []string{"cat", "ls"})

	if tool == nil {
		t.Fatalf("Expected tool to be not nil")
	}

	result1, err4 := tool.Execute("cat", []string{"leak"})

	if result1 != "" || err4 == nil {
		t.Errorf("Expected bare symlink argument to be rejected")
	} else if strings.Contains(err4.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err4)
	} else if strings.Contains(err4.Error(), outside) {
		t.Errorf("Expected error not to leak outside path, got %v", err4)
	} else if strings.Contains(result1, "TOP SECRET") {
		t.Errorf("Expected outside file not to be read, got %s", result1)
	}

	result2, err5 := tool.Execute("cat", []string{"./leak"})

	if result2 != "" || err5 == nil {
		t.Errorf("Expected path symlink argument to be rejected")
	} else if strings.Contains(err5.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err5)
	}

	result3, err6 := tool.Execute("ls", []string{"leakdir"})

	if result3 != "" || err6 == nil {
		t.Errorf("Expected bare symlink directory argument to be rejected")
	} else if strings.Contains(err6.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err6)
	}

	result4, err7 := tool.Execute("ls", []string{"./leakdir"})

	if result4 != "" || err7 == nil {
		t.Errorf("Expected path symlink directory argument to be rejected")
	} else if strings.Contains(err7.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err7)
	}

	result5, err8 := tool.Execute("cat", []string{"normal.txt"})

	if err8 != nil {
		t.Errorf("Expected %v to be nil", err8)
	} else if strings.Contains(result5, "normal content") == false {
		t.Errorf("Expected normal file to be read, got %s", result5)
	}

	result6, err9 := tool.Execute("cat", []string{"./normal.txt"})

	if err9 != nil {
		t.Errorf("Expected %v to be nil", err9)
	} else if strings.Contains(result6, "normal content") == false {
		t.Errorf("Expected normal file to be read, got %s", result6)
	}

	program_tool := NewPrograms([]string{"Execute"}, playground, sandbox, []string{"/bin/cat"})

	if program_tool == nil {
		t.Fatalf("Expected tool to be not nil")
	}

	result7, err10 := program_tool.Execute("/bin/cat", []string{})

	if result7 != "" || err10 == nil {
		t.Errorf("Expected program name with path separator to be rejected")
	} else if strings.Contains(err10.Error(), "Attempt to execute unallowed program") == false {
		t.Errorf("Expected unallowed program error, got %v", err10)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
