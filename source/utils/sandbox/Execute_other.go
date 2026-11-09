//go:build !linux && !darwin

package sandbox

import "fmt"
import "os"

// Execute reports that sandboxing is unsupported on this platform.
func Execute(arguments []string) int {

	fmt.Fprint(os.Stderr, "sandbox.Execute: Sandboxing is not supported on this system.\n")

	return 125

}
