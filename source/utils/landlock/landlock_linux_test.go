//go:build linux

package landlock

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "strings"
import "testing"

func TestMain(main *testing.M) {

	if os.Getenv("EXOCOMP_LANDLOCK_TEST") == "1" {
		os.Exit(runLandlockTestHelper())
	} else {
		os.Exit(main.Run())
	}

}

func runLandlockTestHelper() int {

	sandbox := os.Getenv("EXOCOMP_LANDLOCK_TEST_SANDBOX")
	outside := os.Getenv("EXOCOMP_LANDLOCK_TEST_OUTSIDE")

	runtime.LockOSThread()

	err0 := restrictPaths(DefaultPolicy(sandbox, os.Args[0], nil))

	if err0 != nil {

		fmt.Fprintf(os.Stderr, "restrict: %v\n", err0)
		return 2

	}

	err1 := os.WriteFile(filepath.Join(outside, "pwned.txt"), []byte("pwned"), 0644)

	if err1 == nil {
		fmt.Println("outside-write-ok")
	} else {
		fmt.Println("outside-write-denied")
	}

	_, err2 := os.ReadFile(filepath.Join(outside, "secret.txt"))

	if err2 == nil {
		fmt.Println("outside-read-ok")
	} else {
		fmt.Println("outside-read-denied")
	}

	err3 := os.WriteFile(filepath.Join(sandbox, "inside.txt"), []byte("inside"), 0644)

	if err3 == nil {
		fmt.Println("inside-write-ok")
	} else {
		fmt.Println("inside-write-denied")
	}

	_, err4 := os.ReadFile(filepath.Join(sandbox, "inside.txt"))

	if err4 == nil {
		fmt.Println("inside-read-ok")
	} else {
		fmt.Println("inside-read-denied")
	}

	_, err5 := os.ReadFile(filepath.Join(sandbox, "link-outside"))

	if err5 == nil {
		fmt.Println("symlink-read-ok")
	} else {
		fmt.Println("symlink-read-denied")
	}

	err6 := os.WriteFile(filepath.Join(sandbox, "link-outside"), []byte("pwned"), 0644)

	if err6 == nil {
		fmt.Println("symlink-write-ok")
	} else {
		fmt.Println("symlink-write-denied")
	}

	return 0

}

func TestLandlockRestrictPaths(t *testing.T) {

	if Supported() == false {
		t.Skip("landlock is not available on this kernel")
	}

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-landlock-*")
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

	err3 := os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(sandbox, "link-outside"))

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	cmd := exec.Command(os.Args[0])
	cmd.Env = append(
		os.Environ(),
		"EXOCOMP_LANDLOCK_TEST=1",
		"EXOCOMP_LANDLOCK_TEST_SANDBOX="+sandbox,
		"EXOCOMP_LANDLOCK_TEST_OUTSIDE="+outside,
	)

	output, err4 := cmd.CombinedOutput()

	if err4 != nil {
		t.Fatalf("Expected %v to be nil, got %s", err4, string(output))
	}

	result := string(output)

	if strings.Contains(result, "outside-write-denied") == false {
		t.Errorf("Expected outside writes to be denied, got %s", result)
	}

	if strings.Contains(result, "outside-read-denied") == false {
		t.Errorf("Expected outside reads to be denied, got %s", result)
	}

	if strings.Contains(result, "inside-write-ok") == false {
		t.Errorf("Expected inside writes to be allowed, got %s", result)
	}

	if strings.Contains(result, "inside-read-ok") == false {
		t.Errorf("Expected inside reads to be allowed, got %s", result)
	}

	if strings.Contains(result, "symlink-read-denied") == false {
		t.Errorf("Expected symlink reads to be denied, got %s", result)
	}

	if strings.Contains(result, "symlink-write-denied") == false {
		t.Errorf("Expected symlink writes to be denied, got %s", result)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestDefaultPolicy(t *testing.T) {

	policy := DefaultPolicy("/sandbox", "/usr/bin/cat", []string{"/tmp/scripts"})

	found_sandbox := false
	found_target := false
	found_extra := false

	for _, path := range policy.ReadWrite {
		if path == "/sandbox" {
			found_sandbox = true
		}
	}

	for _, path := range policy.ReadOnly {
		if path == "/usr/bin" {
			found_target = true
		} else if path == "/tmp/scripts" {
			found_extra = true
		}
	}

	if found_sandbox == false {
		t.Errorf("Expected sandbox to be read-write")
	}

	if found_target == false {
		t.Errorf("Expected target folder to be read-only")
	}

	if found_extra == false {
		t.Errorf("Expected extra read path to be read-only")
	}

}

func TestParseArguments(t *testing.T) {

	sandbox, read_paths, target, target_arguments, err0 := parseArguments([]string{
		"--sandbox=/sandbox",
		"--read=/scripts",
		"--",
		"/usr/bin/cat",
		"file.txt",
	})

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if sandbox != "/sandbox" {
		t.Errorf("Expected %q, got %q", "/sandbox", sandbox)
	}

	if len(read_paths) != 1 || read_paths[0] != "/scripts" {
		t.Errorf("Expected read paths %v, got %v", []string{"/scripts"}, read_paths)
	}

	if target != "/usr/bin/cat" {
		t.Errorf("Expected %q, got %q", "/usr/bin/cat", target)
	}

	if len(target_arguments) != 1 || target_arguments[0] != "file.txt" {
		t.Errorf("Expected target arguments %v, got %v", []string{"file.txt"}, target_arguments)
	}

	_, _, _, _, err1 := parseArguments([]string{"--sandbox=/sandbox"})

	if err1 == nil {
		t.Errorf("Expected invalid arguments to fail")
	}

}

func TestAccessFSForABI(t *testing.T) {

	access1 := allAccessFS(1)
	access3 := allAccessFS(3)
	access5 := allAccessFS(5)

	if access1&accessRefer != 0 {
		t.Errorf("Expected ABI 1 not to include REFER")
	}

	if access3&accessTruncate == 0 {
		t.Errorf("Expected ABI 3 to include TRUNCATE")
	}

	if access5&accessIoctlDev == 0 {
		t.Errorf("Expected ABI 5 to include IOCTL_DEV")
	}

}
