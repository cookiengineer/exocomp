package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestResolveSandboxPath(t *testing.T) {

	sandbox, _ := os.MkdirTemp("/tmp", "exocomp-test-resolvesandboxpath-*")
	result1, err1 := resolveSandboxPath(sandbox, "./first")
	result2, err2 := resolveSandboxPath(sandbox, "./first/second")
	result3, err3 := resolveSandboxPath(sandbox, "./")
	result4, err4 := resolveSandboxPath(sandbox, ".")
	result5, err5 := resolveSandboxPath(sandbox, "")
	result6, err6 := resolveSandboxPath(sandbox, sandbox)
	result7, err7 := resolveSandboxPath(sandbox, sandbox+"/first")
	result8, err8 := resolveSandboxPath(sandbox, sandbox+"/first/second")

	if result1 != sandbox+"/first" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result1, sandbox+"/first")
	}

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if result2 != sandbox+"/first/second" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result2, sandbox+"/first/second")
	}

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	if result3 != sandbox {
		t.Errorf("Expected \"%s\" to be \"%s\"", result3, sandbox)
	}

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	if result4 != sandbox {
		t.Errorf("Expected \"%s\" to be \"%s\"", result4, sandbox)
	}

	if err4 != nil {
		t.Errorf("Expected %v to be nil", err4)
	}

	if result5 != sandbox {
		t.Errorf("Expected \"%s\" to be \"%s\"", result5, sandbox)
	}

	if err5 != nil {
		t.Errorf("Expected %v to be nil", err5)
	}

	if result6 != sandbox {
		t.Errorf("Expected \"%s\" to be \"%s\"", result6, sandbox)
	}

	if err6 != nil {
		t.Errorf("Expected %v to be nil", err6)
	}

	if result7 != sandbox+"/first" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result7, sandbox+"/first")
	}

	if err7 != nil {
		t.Errorf("Expected %v to be nil", err7)
	}

	if result8 != sandbox+"/first/second" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result8, sandbox+"/first/second")
	}

	if err8 != nil {
		t.Errorf("Expected %v to be nil", err8)
	}

}

func TestResolveSandboxPathEscapeAttempts(t *testing.T) {

	sandbox, _ := os.MkdirTemp("/tmp", "exocomp-test-resolvesandboxpath-*")
	result1, err1 := resolveSandboxPath(sandbox, "../")
	result2, err2 := resolveSandboxPath(sandbox, "../../first/second")
	result3, err3 := resolveSandboxPath(sandbox, "./../../../../../etc")
	result4, err4 := resolveSandboxPath(sandbox, "\\etc\\secret")

	if result1 != "" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result1, "")
	}

	if err1 == nil {
		t.Errorf("Expected \"%s\" to be not nil", err1)
	} else if err1.Error() != "Invalid path \"..\": Attempt to escape sandbox" {
		t.Errorf("Expected \"%s\" to be \"%s\"", err1.Error(), "Invalid path \"..\": Attempt to escape sandbox")
	}

	if result2 != "" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result2, "")
	}

	if err2 == nil {
		t.Errorf("Expected \"%s\" to be not nil", err2)
	} else if err2.Error() != "Invalid path \"../../first/second\": Attempt to escape sandbox" {
		t.Errorf("Expected \"%s\" to be \"%s\"", err2.Error(), "Invalid path \"../../first/second\": Attempt to escape sandbox")
	}

	if result3 != "" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result3, "")
	}

	if err3 == nil {
		t.Errorf("Expected \"%s\" to be not nil", err3)
	} else if err3.Error() != "Invalid path \"../../etc\": Attempt to escape sandbox" {
		t.Errorf("Expected \"%s\" to be \"%s\"", err3.Error(), "Invalid path \"../../etc\": Attempt to escape sandbox")
	}

	if result4 != "" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result4, "")
	}

	if err4 == nil {
		t.Errorf("Expected \"%s\" to be not nil", err4)
	} else if err4.Error() != "Invalid path \"/etc/secret\": Attempt to escape sandbox" {
		t.Errorf("Expected \"%s\" to be \"%s\"", err4.Error(), "Invalid path \"/etc/secret\": Attempt to escape sandbox")
	}

}

func TestResolveSandboxPathSymlinks(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-resolvesymlinks-*")
	sandbox := filepath.Join(playground, "sandbox")
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

	err3 := os.WriteFile(filepath.Join(sandbox, "internal.txt"), []byte("internal"), 0644)

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	createTestSymlink(t, filepath.Join(outside, "secret.txt"), filepath.Join(sandbox, "link_outside_file"))
	createTestSymlink(t, outside, filepath.Join(sandbox, "link_outside_dir"))
	createTestSymlink(t, "internal.txt", filepath.Join(sandbox, "link_internal"))
	createTestSymlink(t, "missing.txt", filepath.Join(sandbox, "link_dangling"))

	result1, err4 := resolveSandboxPath(sandbox, "link_outside_file")
	result2, err5 := resolveSandboxPath(sandbox, "link_outside_file/child")
	result3, err6 := resolveSandboxPath(sandbox, "link_outside_dir/child")
	result4, err7 := resolveSandboxPath(sandbox, "link_internal")
	result5, err8 := resolveSandboxPath(sandbox, "link_dangling")
	result6, err9 := resolveSandboxPath(sandbox, "nested/new/file.txt")

	if result1 != "" || err4 == nil {
		t.Errorf("Expected outside file symlink to be rejected")
	} else if strings.Contains(err4.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err4)
	} else if strings.Contains(err4.Error(), outside) {
		t.Errorf("Expected error not to leak outside path, got %v", err4)
	}

	if result2 != "" || err5 == nil {
		t.Errorf("Expected child of outside file symlink to be rejected")
	} else if strings.Contains(err5.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err5)
	}

	if result3 != "" || err6 == nil {
		t.Errorf("Expected child of outside dir symlink to be rejected")
	} else if strings.Contains(err6.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err6)
	}

	if result4 != "" || err7 == nil {
		t.Errorf("Expected internal symlink to be rejected")
	} else if strings.Contains(err7.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err7)
	}

	if result5 != "" || err8 == nil {
		t.Errorf("Expected dangling symlink to be rejected")
	} else if strings.Contains(err8.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err8)
	}

	if result6 != filepath.Join(sandbox, "nested", "new", "file.txt") {
		t.Errorf("Expected \"%s\", got \"%s\"", filepath.Join(sandbox, "nested", "new", "file.txt"), result6)
	}

	if err9 != nil {
		t.Errorf("Expected %v to be nil", err9)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
