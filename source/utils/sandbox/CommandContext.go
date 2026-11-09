package sandbox

import "context"
import "fmt"
import "os/exec"
import "path/filepath"

// CommandContext returns an exec.Cmd that runs program through the sandbox
// helper with sandbox as the only writable location.
func CommandContext(ctx context.Context, sandbox string, program string, arguments []string, read_paths []string) (*exec.Cmd, error) {

	target, err0 := exec.LookPath(program)

	if err0 != nil {
		return nil, fmt.Errorf("sandbox.CommandContext: Program doesn't exist.")
	}

	if Supported() == false {
		return nil, fmt.Errorf("sandbox.CommandContext: Sandboxing is not supported on this system.")
	}

	helper, err1 := HelperPath()

	if err1 != nil {
		return nil, err1
	}

	sandbox_abs, err2 := filepath.Abs(sandbox)

	if err2 != nil {
		return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve sandbox path.")
	}

	target_abs, err3 := filepath.Abs(target)

	if err3 != nil {
		return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve program path.")
	}

	environment, err4 := prepareSandboxEnv(sandbox_abs)

	if err4 != nil {
		return nil, err4
	}

	helper_arguments := []string{
		"--sandbox=" + sandbox_abs,
	}

	for _, read_path := range read_paths {

		read_abs, err5 := filepath.Abs(read_path)

		if err5 == nil {
			helper_arguments = append(helper_arguments, "--read="+read_abs)
		} else {
			return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve read path.")
		}

	}

	helper_arguments = append(helper_arguments, "--", target_abs)
	helper_arguments = append(helper_arguments, arguments...)

	cmd := exec.CommandContext(ctx, helper, helper_arguments...)
	cmd.Dir = sandbox_abs
	cmd.Env = environment

	return cmd, nil

}
