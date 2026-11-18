package seatbelt

import "fmt"
import "path/filepath"
import "strings"
import "testing"

func TestBuildProfile(t *testing.T) {

	sandbox := t.TempDir()
	read_only := t.TempDir()

	profile, err0 := BuildProfile(Policy{
		ReadWrite: []string{sandbox},
		ReadOnly:  []string{read_only},
		Devices:   []string{"/dev/null"},
		Deny:      []string{filepath.Join(sandbox, ".exocomp")},
	})

	if err0 == nil {

		source := profile.String()

		checks := []string{
			"(version 1)",
			"(deny default)",
			"(import \"bsd.sb\")",
			"(allow process-fork)",
			"(allow signal)",
			"(allow sysctl-read)",
			"(allow ipc-posix-sem)",
			"com.apple.trustd.agent",
			"com.apple.SystemConfiguration.configd",
			"file-read* file-write* file-read-metadata",
			"file-read* file-read-metadata",
			"process-exec",
			"file-issue-extension",
			"com.apple.app-sandbox.read",
			"com.apple.app-sandbox.read-write",
			"/dev/null",
			"(deny file-read* file-write* file-read-metadata",
		}

		for _, check := range checks {

			if strings.Contains(source, check) == false {
				t.Errorf("Expected profile to contain %q", check)
			}

		}

		denied := []string{
			"network-outbound",
			"network-bind",
			"network-inbound",
			"system-socket",
		}

		for _, check := range denied {

			if strings.Contains(source, check) == true {
				t.Errorf("Expected profile without network access to not contain %q", check)
			}

		}

	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

}

func TestBuildProfileNetwork(t *testing.T) {

	sandbox := t.TempDir()

	profile, err0 := BuildProfile(Policy{
		ReadWrite: []string{sandbox},
		Network:   true,
	})

	if err0 == nil {

		source := profile.String()

		checks := []string{
			"(allow system-socket (socket-domain AF_INET))",
			"(allow system-socket (socket-domain AF_INET6))",
			"(allow network-outbound (remote ip \"*:*\"))",
			"(allow network-outbound (remote unix-socket (path-literal \"/private/var/run/mDNSResponder\")))",
			"(allow network-bind (local ip \"localhost:*\"))",
			"(allow network-inbound (local ip \"localhost:*\"))",
		}

		for _, check := range checks {

			if strings.Contains(source, check) == false {
				t.Errorf("Expected profile to contain %q", check)
			}

		}

		// NOTE: A bare network-outbound rule would also allow unix-socket
		// connects to host daemons such as docker.sock or the SSH agent.
		if strings.Contains(source, "(allow network-outbound)\n") == true {
			t.Errorf("Expected profile to filter network-outbound")
		}

	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

}

func TestBuildProfileInvalidPath(t *testing.T) {

	_, err0 := BuildProfile(Policy{
		ReadWrite: []string{"/tmp/invalid\npath"},
	})

	if err0 != nil {
		// expected: newlines cannot be emitted into SBPL safely
	} else {
		t.Errorf("Expected invalid path to be rejected")
	}

}

func TestBuildProfileTooLarge(t *testing.T) {

	paths := make([]string, 0, 2000)

	for index := 0; index < 2000; index++ {
		paths = append(paths, fmt.Sprintf("/tmp/exocomp-seatbelt-large-%d", index))
	}

	_, err0 := BuildProfile(Policy{
		ReadWrite: paths,
	})

	if err0 != nil {
		// expected: profile exceeds the kernel size limit
	} else {
		t.Errorf("Expected oversized profile to be rejected")
	}

}

func TestQuoteSBPL(t *testing.T) {

	result := quoteSBPL("a\"b\\c")
	expected := "\"a\\\"b\\\\c\""

	if result == expected {
		// expected
	} else {
		t.Errorf("Expected %q, got %q", expected, result)
	}

}

func TestValidatePath(t *testing.T) {

	if validatePath("/tmp/valid") != nil {
		t.Errorf("Expected valid path to be accepted")
	}

	if validatePath("/tmp/invalid\npath") == nil {
		t.Errorf("Expected path with newline to be rejected")
	}

	if validatePath("/tmp/invalid\x00path") == nil {
		t.Errorf("Expected path with NUL to be rejected")
	}

}
