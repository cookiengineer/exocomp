//go:build linux

package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestPrograms_ExecuteLandlock(t *testing.T) {

	requireLandlock(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-landlock-*")
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

	outside_file := filepath.Join(outside, "secret.txt")
	inside_file := filepath.Join(sandbox, "normal.txt")
	outside_target := filepath.Join(outside, "pwned.txt")

	err2 := os.WriteFile(outside_file, []byte("TOP SECRET"), 0644)

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	err3 := os.WriteFile(inside_file, []byte("normal content"), 0644)

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	t.Setenv("EXOCOMP_TEST_OUTSIDE", outside_file)
	t.Setenv("EXOCOMP_TEST_TARGET", outside_target)

	tool := NewPrograms([]string{"Execute", "List", "Stat"}, playground, sandbox, []string{"cat", "sh"})

	if tool == nil {
		t.Fatalf("Expected tool to be not nil")
	}

	result1, err4 := tool.Execute("sh", []string{"-c", "cat \"$EXOCOMP_TEST_OUTSIDE\""})

	if err4 == nil {
		t.Errorf("Expected reading outside the sandbox to fail")
	} else if strings.Contains(result1, "TOP SECRET") {
		t.Errorf("Expected outside file contents not to be returned, got %s", result1)
	}

	_, err5 := tool.Execute("sh", []string{"-c", "echo pwned > \"$EXOCOMP_TEST_TARGET\""})

	if err5 == nil {
		t.Errorf("Expected writing outside the sandbox to fail")
	}

	content1, err6 := os.ReadFile(outside_target)

	if os.IsNotExist(err6) == true {
		// expected: the write was denied
	} else if err6 != nil {
		t.Errorf("Expected %v to be nil", err6)
	} else if strings.Contains(string(content1), "pwned") == true {
		t.Errorf("Expected outside file not to be written")
	}

	_, err7 := tool.Execute("sh", []string{"-c", "echo inside > inside.txt"})

	if err7 != nil {
		t.Errorf("Expected %v to be nil", err7)
	}

	content2, err8 := os.ReadFile(filepath.Join(sandbox, "inside.txt"))

	if err8 != nil {
		t.Errorf("Expected %v to be nil", err8)
	} else if strings.TrimSpace(string(content2)) != "inside" {
		t.Errorf("Expected sandbox file to be written, got %q", string(content2))
	}

	result4, err9 := tool.Execute("cat", []string{"normal.txt"})

	if err9 != nil {
		t.Errorf("Expected %v to be nil", err9)
	} else if strings.Contains(result4, "normal content") == false {
		t.Errorf("Expected sandbox file to be readable, got %s", result4)
	}

	t.Setenv("SSH_AUTH_SOCK", "/tmp/exocomp-fake-agent.sock")

	result5, err10 := tool.Execute("sh", []string{"-c", "printf '%s' \"$SSH_AUTH_SOCK\""})

	if err10 != nil {
		t.Errorf("Expected %v to be nil", err10)
	} else if strings.Contains(result5, "exocomp-fake-agent.sock") == true {
		t.Errorf("Expected agent socket environment to be removed, got %s", result5)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestPrograms_ExecuteLandlockMissingHelper(t *testing.T) {

	requireLandlock(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-landlock-*")
	sandbox := filepath.Join(playground, "programs")

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	err1 := os.WriteFile(filepath.Join(sandbox, "normal.txt"), []byte("normal content"), 0644)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	t.Setenv("EXOCOMP_SANDBOX", filepath.Join(playground, "missing-sandbox"))

	tool := NewPrograms([]string{"Execute", "List", "Stat"}, playground, sandbox, []string{"cat"})

	if tool == nil {
		t.Fatalf("Expected tool to be not nil")
	}

	result1, err2 := tool.Execute("cat", []string{"normal.txt"})

	if result1 != "" {
		t.Errorf("Expected %s to be empty", result1)
	}

	if err2 == nil {
		t.Errorf("Expected %v to be not nil", err2)
	} else if strings.Contains(err2.Error(), "Sandbox helper is not installed.") == false {
		t.Errorf("Expected missing sandbox helper error, got %v", err2)
	} else if strings.Contains(err2.Error(), playground) {
		t.Errorf("Expected error not to leak the host path, got %v", err2)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
