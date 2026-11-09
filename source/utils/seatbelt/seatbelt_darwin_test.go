//go:build darwin

package seatbelt_test

import "exocomp/utils/sandbox"
import "exocomp/utils/seatbelt"
import "context"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "testing"
import "time"

func buildSandboxHelper(t *testing.T) string {

	t.Helper()

	helper := filepath.Join(t.TempDir(), "exocomp-sandbox")
	cmd := exec.Command("go", "build", "-o", helper, "../../cmds/exocomp-sandbox")

	output, err := cmd.CombinedOutput()

	if err == nil {
		return helper
	} else {
		t.Fatalf("Cannot build exocomp-sandbox: %v: %s", err, strings.TrimSpace(string(output)))
		return ""
	}

}

func TestExecuteSandbox(t *testing.T) {

	if seatbelt.Supported() == false {
		t.Skip("seatbelt is not available on this system")
	}

	helper := buildSandboxHelper(t)

	t.Setenv("EXOCOMP_SANDBOX", helper)

	playground := t.TempDir()
	sandbox_dir := filepath.Join(playground, "sandbox")
	outside_dir := filepath.Join(playground, "outside")

	err0 := os.MkdirAll(sandbox_dir, 0755)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	err1 := os.MkdirAll(outside_dir, 0755)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	outside_file := filepath.Join(outside_dir, "secret.txt")
	inside_file := filepath.Join(sandbox_dir, "normal.txt")
	outside_target := filepath.Join(outside_dir, "pwned.txt")

	err2 := os.WriteFile(outside_file, []byte("TOP SECRET"), 0644)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	err3 := os.WriteFile(inside_file, []byte("normal content"), 0644)

	if err3 != nil {
		t.Fatalf("Expected %v to be nil", err3)
	}

	err4 := os.Symlink(outside_file, filepath.Join(sandbox_dir, "link-outside"))

	if err4 != nil {
		t.Fatalf("Expected %v to be nil", err4)
	}

	t.Setenv("EXOCOMP_TEST_OUTSIDE", outside_file)
	t.Setenv("EXOCOMP_TEST_TARGET", outside_target)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	cmd1, err5 := sandbox.CommandContext(ctx, sandbox_dir, "sh", []string{"-c", "cat \"$EXOCOMP_TEST_OUTSIDE\""}, []string{})

	if err5 == nil {

		output, err6 := cmd1.CombinedOutput()

		if err6 == nil {
			t.Errorf("Expected reading outside the sandbox to fail")
		} else if strings.Contains(string(output), "TOP SECRET") == true {
			t.Errorf("Expected outside file contents not to be returned, got %s", string(output))
		}

	} else {
		t.Errorf("Expected %v to be nil", err5)
	}

	cmd2, err7 := sandbox.CommandContext(ctx, sandbox_dir, "sh", []string{"-c", "echo pwned > \"$EXOCOMP_TEST_TARGET\""}, []string{})

	if err7 == nil {

		_, err8 := cmd2.CombinedOutput()

		if err8 == nil {
			t.Errorf("Expected writing outside the sandbox to fail")
		}

	} else {
		t.Errorf("Expected %v to be nil", err7)
	}

	content1, err9 := os.ReadFile(outside_target)

	if os.IsNotExist(err9) == true {
		// expected: the write was denied
	} else if err9 != nil {
		t.Errorf("Expected %v to be nil", err9)
	} else if strings.Contains(string(content1), "pwned") == true {
		t.Errorf("Expected outside file not to be written")
	}

	cmd3, err10 := sandbox.CommandContext(ctx, sandbox_dir, "sh", []string{"-c", "echo inside > inside.txt"}, []string{})

	if err10 == nil {

		_, err11 := cmd3.CombinedOutput()

		if err11 != nil {
			t.Errorf("Expected %v to be nil", err11)
		}

	} else {
		t.Errorf("Expected %v to be nil", err10)
	}

	content2, err12 := os.ReadFile(filepath.Join(sandbox_dir, "inside.txt"))

	if err12 == nil {

		if strings.TrimSpace(string(content2)) != "inside" {
			t.Errorf("Expected sandbox file to be written, got %q", string(content2))
		}

	} else {
		t.Errorf("Expected %v to be nil", err12)
	}

	cmd4, err13 := sandbox.CommandContext(ctx, sandbox_dir, "cat", []string{"normal.txt"}, []string{})

	if err13 == nil {

		output, err14 := cmd4.CombinedOutput()

		if err14 != nil {
			t.Errorf("Expected %v to be nil", err14)
		} else if strings.Contains(string(output), "normal content") == false {
			t.Errorf("Expected sandbox file to be readable, got %s", string(output))
		}

	} else {
		t.Errorf("Expected %v to be nil", err13)
	}

	cmd5, err15 := sandbox.CommandContext(ctx, sandbox_dir, "cat", []string{"link-outside"}, []string{})

	if err15 == nil {

		output, err16 := cmd5.CombinedOutput()

		if err16 == nil {
			t.Errorf("Expected reading through symlink to fail")
		} else if strings.Contains(string(output), "TOP SECRET") == true {
			t.Errorf("Expected symlink target contents not to be returned, got %s", string(output))
		}

	} else {
		t.Errorf("Expected %v to be nil", err15)
	}

}
