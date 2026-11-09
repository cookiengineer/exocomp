//go:build !linux

package landlock

import "fmt"
import "os"

// Execute reports that landlock sandboxing is unsupported on this platform.
func Execute(arguments []string) int {

	fmt.Fprint(os.Stderr, "landlock.Execute: Sandboxing is not supported on this system.\n")

	return 125

}
