package sandbox

import "context"
import "os/exec"

// CommandContext returns an exec.Cmd that runs program through the sandbox
// helper with sandbox as the only writable location.
func CommandContext(ctx context.Context, sandbox string, program string, arguments []string, read_paths []string) (*exec.Cmd, error) {

	return CommandContextWithPermissions(ctx, sandbox, program, arguments, Permissions{
		ReadPaths: read_paths,
	})

}
