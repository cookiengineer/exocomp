package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestPrograms_Execute(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox       := filepath.Join(playground, "programs")
	tool          := NewPrograms("tester", sandbox, []string{"cat", "ls", "pwd"})

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if tool != nil {

		result1, err1 := tool.Execute("ls", []string{"-la"})
		result2, err2 := tool.Execute("pwd", []string{})
		result3, err3 := tool.Execute("cat", []string{"./../../../file.txt"})
		result4, err4 := tool.Execute("cat", []string{"/etc/passwd"})
		result5, err5 := tool.Execute("nc", []string{"192.168.123.123", "1337"})

		lines1 := strings.Split(result1, "\n")

		if len(lines1) == 5 {

			if lines1[0] != "programs.Execute: ls -la" {
				t.Errorf("Expected \"%s\" to be executed", "ls -la")
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

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-*")
	sandbox       := filepath.Join(playground, "programs")
	tool          := NewPrograms("tester", sandbox, []string{"ls", "pwd"})

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
	sandbox       := filepath.Join(playground, "programs")
	tool          := NewPrograms("tester", sandbox, []string{"doesntexist"})

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

