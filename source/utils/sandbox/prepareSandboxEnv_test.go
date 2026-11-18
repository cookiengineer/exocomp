package sandbox

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestPrepareSandboxEnv_RemovesSecrets(t *testing.T) {

	playground, err0 := os.MkdirTemp("/tmp", "exocomp-test-sandbox-env-*")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	t.Cleanup(func() {
		os.RemoveAll(playground)
	})

	secrets := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"PDCP_API_KEY",
		"SLACK_WEBHOOK_URL",
	}

	for _, name := range secrets {
		t.Setenv(name, "secret")
	}

	t.Setenv("EXOCOMP_TEST_VISIBLE", "visible")

	environment, _, err1 := prepareSandboxEnv(playground)

	if err1 == nil {

		for _, name := range secrets {

			if environmentContains(environment, name) == true {
				t.Errorf("Expected environment variable %q to be removed", name)
			}

		}

		if environmentContains(environment, "EXOCOMP_TEST_VISIBLE") == false {
			t.Errorf("Expected unrelated environment variable to be preserved")
		}

	} else {
		t.Errorf("Expected %v to be nil", err1)
	}

}

func TestPrepareSandboxEnv_UsesRuntimeOutsideSandbox(t *testing.T) {

	playground, err0 := os.MkdirTemp("/tmp", "exocomp-test-sandbox-runtime-*")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	t.Cleanup(func() {
		os.RemoveAll(playground)
	})

	_, runtime_paths, err1 := prepareSandboxEnv(playground)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if len(runtime_paths) != 1 {
		t.Fatalf("Expected one runtime path, got %v", runtime_paths)
	}

	stat, err2 := os.Stat(runtime_paths[0])

	if err2 != nil {
		t.Fatalf("Expected runtime directory to exist: %v", err2)
	}

	if stat.IsDir() == false {
		t.Errorf("Expected runtime path to be a directory")
	}

	if strings.HasPrefix(runtime_paths[0], playground) == true {
		t.Errorf("Expected runtime directory to be outside the sandbox")
	}

	if _, err3 := os.Stat(filepath.Join(playground, ".exocomp")); os.IsNotExist(err3) == false {
		t.Errorf("Expected prepareSandboxEnv to not create .exocomp")
	}

}

func environmentContains(environment []string, name string) bool {

	prefix := name + "="

	for _, entry := range environment {

		if strings.HasPrefix(entry, prefix) == true {
			return true
		}

	}

	return false

}
