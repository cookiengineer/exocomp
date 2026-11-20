package testutil

import "exocomp/utils/sandbox"
import "testing"

// RequireSandbox skips the test when the platform sandbox is unavailable and
// otherwise points EXOCOMP_SANDBOX at a freshly built helper.
func RequireSandbox(t *testing.T) {

	t.Helper()

	if sandbox.Supported() == false {
		t.Skip("sandboxing is not available on this system")
	}

	t.Setenv("EXOCOMP_SANDBOX", BuildSandboxExec(t))

}
