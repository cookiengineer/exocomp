package sandbox

import "context"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"

// CommandContextWithPermissions returns an exec.Cmd that runs program through
// the sandbox helper with extra filesystem permissions and, when Network is
// true, network access. Network access is denied by default.
func CommandContextWithPermissions(ctx context.Context, sandbox string, program string, arguments []string, permissions Permissions) (*exec.Cmd, error) {

	target, err0 := exec.LookPath(program)

	if err0 == nil {

		if Supported() == true {

			helper, err1 := HelperPath()

			if err1 == nil {

				sandbox_abs, err2 := filepath.Abs(sandbox)

				if err2 != nil {
					return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve sandbox path.")
				}

				err_mkdir := os.MkdirAll(sandbox_abs, 0755)

				if err_mkdir != nil {
					return nil, fmt.Errorf("sandbox.CommandContext: Permission denied.")
				}

				stat, err_stat := os.Stat(sandbox_abs)

				if err_stat != nil || stat.IsDir() == false || stat.Mode().Perm()&0111 == 0 {
					return nil, fmt.Errorf("sandbox.CommandContext: Permission denied.")
				}

				target_abs, err3 := filepath.Abs(target)

				if err3 != nil {
					return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve program path.")
				}

				environment, runtime_paths, err4 := prepareSandboxEnv(sandbox_abs)

				if err4 != nil {
					return nil, err4
				}

				helper_arguments := []string{
					"--sandbox=" + sandbox_abs,
				}

				for _, read_path := range permissions.ReadPaths {

					read_abs, err5 := filepath.Abs(read_path)

					if err5 == nil {
						helper_arguments = append(helper_arguments, "--read="+read_abs)
					} else {
						return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve read path.")
					}

				}

				write_paths := append([]string{}, permissions.WritePaths...)
				write_paths = append(write_paths, runtime_paths...)

				for _, write_path := range write_paths {

					write_abs, err6 := filepath.Abs(write_path)

					if err6 == nil {
						helper_arguments = append(helper_arguments, "--write="+write_abs)
					} else {
						return nil, fmt.Errorf("sandbox.CommandContext: Cannot resolve write path.")
					}

				}

				if permissions.Network == true {
					helper_arguments = append(helper_arguments, "--network=allow")
				} else {
					helper_arguments = append(helper_arguments, "--network=deny")
				}

				helper_arguments = append(helper_arguments, "--", target_abs)
				helper_arguments = append(helper_arguments, arguments...)

				cmd := exec.CommandContext(ctx, helper, helper_arguments...)
				cmd.Dir = sandbox_abs
				cmd.Env = environment

				prepareSandboxCommand(cmd, permissions.Network)

				return cmd, nil

			} else {
				return nil, err1
			}

		} else {
			return nil, fmt.Errorf("sandbox.CommandContext: Sandboxing is not supported on this system.")
		}
	} else {
		return nil, fmt.Errorf("sandbox.CommandContext: Program doesn't exist.")
	}

}
