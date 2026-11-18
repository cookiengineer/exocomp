package sandbox

import "os"
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

	environment, err1 := prepareSandboxEnv(playground)

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

func environmentContains(environment []string, name string) bool {

	prefix := name + "="

	for _, entry := range environment {

		if strings.HasPrefix(entry, prefix) == true {
			return true
		}

	}

	return false

}
