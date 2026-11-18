//go:build linux

package sandbox

import "context"
import "errors"
import "fmt"
import "net"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "testing"
import "time"

// TestSandboxNetworkHelper is re-executed as the sandbox target by
// TestSandboxNetwork. It prints a marker so the parent can tell whether the
// sandboxed process can reach the host or its own loopback interface.
func TestSandboxNetworkHelper(t *testing.T) {

	mode := os.Getenv("EXOCOMP_NETWORK_TEST_MODE")

	if mode == "" {
		t.Skip("not a sandbox network helper invocation")
	}

	if mode == "dial" {

		conn, err0 := net.DialTimeout("tcp", os.Getenv("EXOCOMP_NETWORK_TEST_ADDRESS"), 2*time.Second)

		if err0 == nil {
			fmt.Println("dial-ok")
			conn.Close()
		} else {
			fmt.Println("dial-denied")
		}

	} else if mode == "loopback" {

		listener, err0 := net.Listen("tcp", "127.0.0.1:0")

		if err0 != nil {
			fmt.Println("loopback-listen-denied")
			return
		}

		defer listener.Close()

		go func() {

			for {

				connection, err := listener.Accept()

				if err != nil {
					return
				}

				connection.Close()

			}

		}()

		conn, err1 := net.DialTimeout("tcp", listener.Addr().String(), 2*time.Second)

		if err1 == nil {
			fmt.Println("loopback-ok")
			conn.Close()
		} else {
			fmt.Println("loopback-denied")
		}

	} else {
		fmt.Println("unknown-mode")
	}

}

func buildNetworkTestHelper(t *testing.T) string {

	t.Helper()

	helper := filepath.Join(t.TempDir(), "exocomp-sandbox")
	cmd := exec.Command("go", "build", "-o", helper, "../../cmds/exocomp-sandbox")

	output, err0 := cmd.CombinedOutput()

	if err0 == nil {
		return helper
	} else {
		t.Fatalf("Cannot build exocomp-sandbox: %v: %s", err0, strings.TrimSpace(string(output)))
		return ""
	}

}

func runNetworkTestHelper(t *testing.T, sandbox string, mode string, address string, network bool) (string, error) {

	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Setenv("EXOCOMP_NETWORK_TEST_MODE", mode)
	t.Setenv("EXOCOMP_NETWORK_TEST_ADDRESS", address)

	cmd, err0 := CommandContextWithPermissions(
		ctx,
		sandbox,
		os.Args[0],
		[]string{"-test.run=^TestSandboxNetworkHelper$"},
		Permissions{Network: network},
	)

	if err0 == nil {

		output, err1 := cmd.CombinedOutput()

		return string(output), err1

	} else {
		return "", err0
	}

}

func TestSandboxNetwork(t *testing.T) {

	if Supported() == false {
		t.Skip("landlock is not available on this kernel")
	}

	helper := buildNetworkTestHelper(t)

	t.Setenv("EXOCOMP_SANDBOX", helper)

	sandbox := filepath.Join(t.TempDir(), "sandbox")

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	probe_ctx, probe_cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer probe_cancel()

	probe, err_probe := CommandContextWithPermissions(
		probe_ctx,
		sandbox,
		"true",
		[]string{},
		Permissions{},
	)

	if err_probe != nil {
		t.Skipf("Sandboxing is not available on this system: %v", err_probe)
	}

	probe_output, err_probe_run := probe.CombinedOutput()
	probe_text := string(probe_output)

	if err_probe_run != nil {

		var exit_error *exec.ExitError

		if errors.As(err_probe_run, &exit_error) == false {
			t.Skipf("User and network namespaces are not available: %v", err_probe_run)
		}

		if strings.Contains(probe_text, "Sandboxing is not supported") == true ||
			strings.Contains(probe_text, "Cannot isolate network namespace") == true ||
			strings.Contains(probe_text, "Cannot start program") == true {
			t.Skipf("User and network namespaces are not available: %s", strings.TrimSpace(probe_text))
		}

		t.Fatalf("Sandbox probe failed: %v: %s", err_probe_run, strings.TrimSpace(probe_text))

	}

	listener, err1 := net.Listen("tcp", "127.0.0.1:0")

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	defer listener.Close()

	go func() {

		for {

			connection, err := listener.Accept()

			if err != nil {
				return
			}

			connection.Close()

		}

	}()

	address := listener.Addr().String()

	output_denied, err_denied := runNetworkTestHelper(t, sandbox, "dial", address, false)

	if strings.Contains(output_denied, "dial-denied") == false {
		t.Errorf("Expected network access to be denied, got %s (%v)", output_denied, err_denied)
	}

	output_loopback, err_loopback := runNetworkTestHelper(t, sandbox, "loopback", "", false)

	if strings.Contains(output_loopback, "loopback-ok") == false {
		t.Errorf("Expected loopback to work in isolated mode, got %s (%v)", output_loopback, err_loopback)
	}

	output_allowed, err_allowed := runNetworkTestHelper(t, sandbox, "dial", address, true)

	if strings.Contains(output_allowed, "dial-ok") == false {
		t.Errorf("Expected network access to be allowed, got %s (%v)", output_allowed, err_allowed)
	}

}
