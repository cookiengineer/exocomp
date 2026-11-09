package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestSanitizeSandboxPath(t *testing.T) {

	sandbox, _ := os.MkdirTemp("/tmp", "exocomp-test-sanitizesandboxpath-*")
	result1, err1 := sanitizeSandboxPath(sandbox, "./first")
	result2, err2 := sanitizeSandboxPath(sandbox, "./first/second")
	result3, err3 := sanitizeSandboxPath(sandbox, "./")
	result4, err4 := sanitizeSandboxPath(sandbox, ".")
	result5, err5 := sanitizeSandboxPath(sandbox, "")
	result6, err6 := sanitizeSandboxPath(sandbox, sandbox)
	result7, err7 := sanitizeSandboxPath(sandbox, sandbox+"/first")
	result8, err8 := sanitizeSandboxPath(sandbox, sandbox+"/first/second")

	if result1 != "./first" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result1, "./first")
	}

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if result2 != "./first/second" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result2, "./first/second")
	}

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	if result3 != "." {
		t.Errorf("Expected \"%s\" to be \"%s\"", result3, ".")
	}

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	if result4 != "." {
		t.Errorf("Expected \"%s\" to be \"%s\"", result4, ".")
	}

	if err4 != nil {
		t.Errorf("Expected %v to be nil", err4)
	}

	if result5 != "." {
		t.Errorf("Expected \"%s\" to be \"%s\"", result5, ".")
	}

	if err5 != nil {
		t.Errorf("Expected %v to be nil", err5)
	}

	if result6 != "." {
		t.Errorf("Expected \"%s\" to be \"%s\"", result6, ".")
	}

	if err6 != nil {
		t.Errorf("Expected %v to be nil", err6)
	}

	if result7 != "./first" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result7, "./first")
	}

	if err7 != nil {
		t.Errorf("Expected %v to be nil", err7)
	}

	if result8 != "./first/second" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result8, "./first/second")
	}

	if err8 != nil {
		t.Errorf("Expected %v to be nil", err8)
	}

}

func TestSanitizeSandboxPathEscapeAttempts(t *testing.T) {

	sandbox, _ := os.MkdirTemp("/tmp", "exocomp-test-sanitizesandboxpath-*")
	result1, err1 := sanitizeSandboxPath(sandbox, "../")
	result2, err2 := sanitizeSandboxPath(sandbox, "../../first/second")
	result3, err3 := sanitizeSandboxPath(sandbox, "./../../../../../etc")
	result4, err4 := sanitizeSandboxPath(sandbox, "\\etc\\secret")

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

func TestSanitizeSandboxPathSymlinks(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-sanitizesymlinks-*")
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

	createTestSymlink(t, outside, filepath.Join(sandbox, "link_outside_dir"))
	createTestSymlink(t, outside, filepath.Join(sandbox, "link_nested_dir"))

	result1, err2 := sanitizeSandboxPath(sandbox, "link_outside_dir/child")
	result2, err3 := sanitizeSandboxPath(sandbox, "link_nested_dir/child")
	result3, err4 := sanitizeSandboxPath(sandbox, "newdir/newfile.txt")

	if result1 != "" || err2 == nil {
		t.Errorf("Expected outside dir symlink to be rejected")
	} else if strings.Contains(err2.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err2)
	}

	if result2 != "" || err3 == nil {
		t.Errorf("Expected outside dir symlink to be rejected")
	} else if strings.Contains(err3.Error(), "Symbolic links are not allowed") == false {
		t.Errorf("Expected symlink error, got %v", err3)
	}

	if result3 != "./newdir/newfile.txt" {
		t.Errorf("Expected \"%s\", got \"%s\"", "./newdir/newfile.txt", result3)
	}

	if err4 != nil {
		t.Errorf("Expected %v to be nil", err4)
	}

	_, err5 := os.Stat(filepath.Join(sandbox, "newdir"))

	if os.IsNotExist(err5) == false {
		t.Errorf("Expected sanitizeSandboxPath not to create folders")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
