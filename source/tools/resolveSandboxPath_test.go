package tools

import "os"
import "testing"

func TestResolveSandboxPath(t *testing.T) {

	sandbox, _    := os.MkdirTemp("/tmp", "exocomp-test-resolvesandboxpath-*")
	result1, err1 := resolveSandboxPath(sandbox, "./first")
	result2, err2 := resolveSandboxPath(sandbox, "./first/second")
	result3, err3 := resolveSandboxPath(sandbox, "./")
	result4, err4 := resolveSandboxPath(sandbox, ".")
	result5, err5 := resolveSandboxPath(sandbox, "")
	result6, err6 := resolveSandboxPath(sandbox, sandbox)
	result7, err7 := resolveSandboxPath(sandbox, sandbox + "/first")
	result8, err8 := resolveSandboxPath(sandbox, sandbox + "/first/second")

	if result1 != sandbox + "/first" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result1, sandbox + "/first")
	}

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if result2 != sandbox + "/first/second" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result2, sandbox + "/first/second")
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

	if result7 != sandbox + "/first" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result7, sandbox + "/first")
	}

	if err7 != nil {
		t.Errorf("Expected %v to be nil", err7)
	}

	if result8 != sandbox + "/first/second" {
		t.Errorf("Expected \"%s\" to be \"%s\"", result8, sandbox + "/first/second")
	}

	if err8 != nil {
		t.Errorf("Expected %v to be nil", err8)
	}

}

func TestResolveSandboxPathEscapeAttempts(t *testing.T) {

	sandbox, _    := os.MkdirTemp("/tmp", "exocomp-test-resolvesandboxpath-*")
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
