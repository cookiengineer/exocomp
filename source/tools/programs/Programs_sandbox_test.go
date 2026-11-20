package programs

import "exocomp/tools/testutil"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestPrograms_ExecuteSandbox(t *testing.T) {

	testutil.RequireSandbox(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-sandbox-*")
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

func TestPrograms_ExecuteHidesExocomp(t *testing.T) {

	testutil.RequireSandbox(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-exocomp-*")
	sandbox := filepath.Join(playground, "programs")
	private := filepath.Join(sandbox, ".exocomp")

	if err := os.MkdirAll(private, 0755); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	secret := filepath.Join(private, "secret.txt")

	if err := os.WriteFile(secret, []byte("TOP SECRET"), 0644); err != nil {
		t.Fatalf("Expected %v to be nil", err)
	}

	tool := NewPrograms([]string{"Execute"}, playground, sandbox, []string{"cat", "sh"})

	result1, err1 := tool.Execute("sh", []string{"-c", "test -e .exocomp/secret.txt && echo visible || echo hidden"})

	if err1 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err1, result1)
	}

	if strings.Contains(result1, "hidden") == false {
		t.Errorf("Expected .exocomp to be hidden, got %s", result1)
	}

	result2, err2 := tool.Execute("cat", []string{".exocomp/secret.txt"})

	if err2 == nil {
		t.Errorf("Expected reading .exocomp to fail")
	} else if strings.Contains(result2, "TOP SECRET") == true {
		t.Errorf("Expected secret content not to be returned, got %s", result2)
	}

	t.Setenv("EXOCOMP_TEST_PROC_ROOT", "/proc/1/root"+secret)

	result_proc, err_proc := tool.Execute("sh", []string{"-c", "cat \"$EXOCOMP_TEST_PROC_ROOT\""})

	if err_proc == nil {
		t.Errorf("Expected reading .exocomp through /proc to fail")
	} else if strings.Contains(result_proc, "TOP SECRET") == true {
		t.Errorf("Expected secret content not to be returned through /proc, got %s", result_proc)
	}

	_, err3 := tool.Execute("sh", []string{"-c", "echo pwned > .exocomp/pwned.txt"})

	if err3 == nil {
		t.Errorf("Expected writing into .exocomp to fail")
	}

	if _, err4 := os.Stat(filepath.Join(private, "pwned.txt")); os.IsNotExist(err4) == false {
		t.Errorf("Expected host .exocomp not to be written")
	}

	result5, err5 := tool.Execute("sh", []string{"-c", "cd \"$HOME\" && echo runtime > runtime.txt && cat runtime.txt"})

	if err5 != nil {
		t.Fatalf("Expected runtime HOME to be writable, got %s", result5)
	}

	if strings.Contains(result5, "runtime") == false {
		t.Errorf("Expected runtime HOME content, got %s", result5)
	}

	if content, err6 := os.ReadFile(secret); err6 != nil || strings.Contains(string(content), "TOP SECRET") == false {
		t.Errorf("Expected host .exocomp to stay intact")
	}

	t.Cleanup(func() {
		os.RemoveAll(playground)
	})

}

func TestPrograms_ExecuteSandboxMissingHelper(t *testing.T) {

	testutil.RequireSandbox(t)

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-programs-sandbox-*")
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
