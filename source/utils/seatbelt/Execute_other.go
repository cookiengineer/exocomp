//go:build !darwin

package seatbelt

import "fmt"
import "os"

// Execute reports that seatbelt sandboxing is unsupported on this platform.
func Execute(arguments []string) int {

	fmt.Fprint(os.Stderr, "seatbelt.Execute: Sandboxing is not supported on this system.\n")

	return 125

}
